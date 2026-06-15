package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"web-chat/internal/domain"
	"web-chat/internal/hub"
	"web-chat/internal/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(h *hub.Hub, msgSvc service.MessageService, userSvc service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("websocket upgrade failed", "error", err)
			return
		}

		roomIDStr := r.PathValue("room_id")
		userIDStr := r.URL.Query().Get("user_id")

		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			slog.Error("invalid room_id", "value", roomIDStr)
			conn.Close()
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			slog.Error("invalid user_id", "value", userIDStr)
			conn.Close()
			return
		}

		client := &hub.Client{
			Conn:   conn,
			UserID: userID,
			RoomID: roomID,
			Send:   make(chan []byte, 256),
		}

		h.Register <- client

		go writePump(client)
		go readPump(client, h, msgSvc, userSvc)
	}
}

func readPump(client *hub.Client, h *hub.Hub, msgSvc service.MessageService, userSvc service.UserService) {
	defer func() {
		h.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var incoming struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}

		if err := json.Unmarshal(message, &incoming); err != nil {
			slog.Error("invalid message format", "error", err)
			continue
		}

		if incoming.Type != "message" {
			continue
		}

		msg := &domain.Message{
			Text:   incoming.Text,
			RoomID: client.RoomID,
			UserID: client.UserID,
		}

		if err := msgSvc.Create(context.Background(), msg); err != nil {
			slog.Error("failed to save message", "error", err)
			continue
		}

		username := "unknown"
		user, err := userSvc.GetByID(context.Background(), client.UserID)
		if err != nil {
			slog.Warn("failed to get user for broadcast", "user_id", client.UserID)
		}

		if user != nil {
			username = user.FullName
		}

		response := map[string]any{
			"type":       "message",
			"id":         msg.ID,
			"user_id":    client.UserID,
			"user":       username,
			"text":       msg.Text,
			"created_at": msg.CreatedAt,
		}

		data, _ := json.Marshal(response)
		h.Broadcast <- hub.BroadcastMsg{
			RoomID: msg.RoomID,
			Data:   data,
		}
	}
}

func writePump(client *hub.Client) {
	defer client.Conn.Close()

	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
