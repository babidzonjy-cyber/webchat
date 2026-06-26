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

type ClientSet map[*Client]struct{}

type Hub struct {
	RoomClients      map[int]ClientSet
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
		RoomClients:      make(map[int]ClientSet),
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
			if err := h.register(client); err != nil {
				slog.Error("redis add online", "room", client.RoomID, "user", client.UserID, "error", err)
			}
		case client := <-h.Unregister:
			if err := h.unregister(client); err != nil {
				slog.Error("redis remove online", "room", client.RoomID, "user", client.UserID, "error", err)
			}
		case msg := <-h.Broadcast:
			h.broadcast(msg)
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

func (h *Hub) register(client *Client) error {
	if _, exists := h.RoomClients[client.RoomID]; !exists {
		h.RoomClients[client.RoomID] = make(map[*Client]struct{})
	}

	h.RoomClients[client.RoomID][client] = struct{}{}

	return h.online.AddOnline(client.RoomID, client.UserID)
}

func (h *Hub) unregister(client *Client) error {
	if _, exists := h.RoomClients[client.RoomID][client]; exists {
		delete(h.RoomClients[client.RoomID], client)
		close(client.Send)

		if len(h.RoomClients[client.RoomID]) == 0 {
			delete(h.RoomClients, client.RoomID)
		}
	}
	return h.online.RemoveOnline(client.RoomID, client.UserID)
}

func (h *Hub) broadcast(msg BroadcastMsg) {
	for client := range h.RoomClients[msg.RoomID] {
		select {
		case client.Send <- msg.Data:
		default:
			if err := h.unregister(client); err != nil {
				slog.Error("redis remove online", "room", client.RoomID, "user", client.UserID, "error", err)
			}
		}
	}
}
