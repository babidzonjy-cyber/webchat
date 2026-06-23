package hub

import (
	"log/slog"
	"web-chat/internal/repository"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID int
	RoomID int
	Send   chan []byte
}

type OnlineCountRequest struct {
	RoomID int
	Result chan int
}

type RoomClientsRequest struct {
	RoomID int
	Result chan []int
}

type IsUserOnlineRequest struct {
	UserID int
	RoomID int
	Result chan bool
}

type Hub struct {
	Clients          map[*Client]bool
	Broadcast        chan BroadcastMsg
	Register         chan *Client
	Unregister       chan *Client
	GetCount         chan OnlineCountRequest
	GetClientsInRoom chan RoomClientsRequest
	CheckUserOnline  chan IsUserOnlineRequest

	online repository.OnlineRepository
}

type BroadcastMsg struct {
	RoomID int
	Data   []byte
}

func NewHub(online repository.OnlineRepository) *Hub {
	return &Hub{
		Clients:          make(map[*Client]bool),
		Broadcast:        make(chan BroadcastMsg, 256),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		GetCount:         make(chan OnlineCountRequest),
		GetClientsInRoom: make(chan RoomClientsRequest),
		CheckUserOnline:  make(chan IsUserOnlineRequest),
		online:           online,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			if err := h.online.AddOnline(client.RoomID, client.UserID); err != nil {
				slog.Error("redis add online", "room", client.RoomID, "user", client.UserID, "error", err)
			}
		case client := <-h.Unregister:
			if _, exists := h.Clients[client]; exists {
				if err := h.online.RemoveOnline(client.RoomID, client.UserID); err != nil {
					slog.Error("redis remove online", "room", client.RoomID, "user", client.UserID, "error", err)
				}
				delete(h.Clients, client)
				close(client.Send)
			}
		case msg := <-h.Broadcast:
			for client := range h.Clients {
				if client.RoomID == msg.RoomID {
					select {
					case client.Send <- msg.Data:
					default:
						close(client.Send)
						delete(h.Clients, client)
						if err := h.online.RemoveOnline(client.RoomID, client.UserID); err != nil {
							slog.Error("broadcast redis remove online", "room", client.RoomID, "user", client.UserID, "error", err)
						}
					}
				}
			}
		case req := <-h.GetCount:
			count, err := h.online.GetOnlineCount(req.RoomID)
			if err != nil {
				slog.Error("redis get online count", "room", req.RoomID, "error", err)
			}
			req.Result <- count
		case req := <-h.GetClientsInRoom:
			ids, err := h.online.GetOnlineUsers(req.RoomID)
			if err != nil {
				slog.Error("redis get online users", "room", req.RoomID, "error", err)
			}
			req.Result <- ids
		case req := <-h.CheckUserOnline:
			online, err := h.online.IsOnline(req.UserID, req.RoomID)
			if err != nil {
				slog.Error("redis get online users", "room", req.RoomID, "user", req.UserID, "error", err)
			}
			req.Result <- online
		}
	}
}

func (h *Hub) GetOnlineCount(roomID int) int {
	req := OnlineCountRequest{
		RoomID: roomID,
		Result: make(chan int),
	}
	h.GetCount <- req
	return <-req.Result
}

func (h *Hub) CheckOnline(userID, roomID int) bool {
	req := IsUserOnlineRequest{
		UserID: userID,
		RoomID: roomID,
		Result: make(chan bool),
	}
	h.CheckUserOnline <- req
	return <-req.Result
}

func (h *Hub) GetUsersInRoom(roomID int) []int {
	req := RoomClientsRequest{
		RoomID: roomID,
		Result: make(chan []int),
	}
	h.GetClientsInRoom <- req
	return <-req.Result
}
