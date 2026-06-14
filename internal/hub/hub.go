package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID int
	RoomID int
	Send   chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan BroadcastMsg
	Register   chan *Client
	Unregister chan *Client

	mtx sync.RWMutex
}

type BroadcastMsg struct {
	RoomID int
	Data   []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan BroadcastMsg, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
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
		}
	}
}

func (h *Hub) GetOnlineCount(roomID int) int {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	count := 0
	for client := range h.Clients {
		if client.RoomID == roomID {
			count++
		}
	}

	return count
}
