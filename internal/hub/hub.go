package hub

import (
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
}

type BroadcastMsg struct {
	RoomID int
	Data   []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:          make(map[*Client]bool),
		Broadcast:        make(chan BroadcastMsg, 256),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		GetCount:         make(chan OnlineCountRequest),
		GetClientsInRoom: make(chan RoomClientsRequest),
		CheckUserOnline:  make(chan IsUserOnlineRequest),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, exists := h.Clients[client]; exists {
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
					}
				}
			}
		case req := <-h.GetCount:
			count := 0
			for client := range h.Clients {
				if client.RoomID == req.RoomID {
					count++
				}
			}
			req.Result <- count
		case req := <-h.GetClientsInRoom:
			clientsID := make([]int, 0)
			for client := range h.Clients {
				if client.RoomID == req.RoomID {
					clientsID = append(clientsID, client.UserID)
				}
			}
			req.Result <- clientsID
		case req := <-h.CheckUserOnline:
			online := false
			for client := range h.Clients {
				if client.UserID == req.UserID && client.RoomID == req.RoomID {
					online = true
					break
				}
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
