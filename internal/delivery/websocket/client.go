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

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func readPump(client *hub.Client, h *hub.Hub, msgSvc service.MessageService, userSvc service.UserService, pool *worker.Pool) {
	defer func() {
		h.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

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
				slog.Error("failed to proccess message", "error", err)
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.TextMessage, []byte{})
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				slog.Error(
					"failed to write websocket message",
					"error", err,
					"user_id", client.UserID,
					"room_id", client.RoomID,
				)
				return
			}
		case <-ticker.C:
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
