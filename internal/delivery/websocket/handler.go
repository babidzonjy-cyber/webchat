package websocket

import (
	"log/slog"
	"net/http"
	"strconv"
	"web-chat/internal/auth"
	"web-chat/internal/hub"
	"web-chat/internal/service"
	"web-chat/internal/worker"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(h *hub.Hub, msgSvc service.MessageService, userSvc service.UserService, pool *worker.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("websocket upgrade failed", "error", err)
			return
		}

		roomIDStr := r.PathValue("room_id")
		userID := auth.UserIDFromContext(r.Context())

		if userID == 0 {
			conn.Close()
			return
		}

		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			slog.Error("invalid room_id", "value", roomIDStr)
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
		go readPump(client, h, msgSvc, userSvc, pool)
	}
}
