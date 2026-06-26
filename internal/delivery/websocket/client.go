package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
	"web-chat/internal/hub"
	"web-chat/internal/service"
	"web-chat/internal/worker"

	"github.com/gorilla/websocket"
)

func readPump(client *hub.Client, h *hub.Hub, msgSvc service.MessageService, userSvc service.UserService, pool *worker.Pool) {
	defer func() {
		h.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(message, &incoming); err != nil {
			slog.Error("invalid message format", "error", err)
			continue
		}

		if incoming.Type != "message" {
			continue
		}

		pool.Submit(func() {
			msg := buildDomainMessage(client, incoming)

			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()

			response, err := handleIncomingMessage(ctx, msgSvc, userSvc, msg, client)
			if err != nil {
				slog.Error("error", err)
				errMsg := ErrorMessage{
					Type:    "error",
					Message: "failed to process message",
				}

				errData, _ := json.Marshal(errMsg)
				client.Send <- errData
				return
			}

			data, _ := json.Marshal(response)
			h.Broadcast <- hub.BroadcastMsg{
				RoomID: msg.RoomID,
				Data:   data,
			}
		})
	}
}

func writePump(client *hub.Client) {
	defer client.Conn.Close()

	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			slog.Error(
				"failed to write websocket message",
				"error", err,
				"user_id", client.UserID,
				"room_id", client.RoomID,
			)
			return
		}
	}
}
