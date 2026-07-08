package websocket

import (
	"log/slog"
	"net/http"
	"strconv"
	"web-chat/internal/auth"
	"web-chat/internal/domain"
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

func ServeWS(
	h *hub.Hub,
	msgSvc service.MessageService,
	userSvc service.UserService,
	pool *worker.Pool,
	roomMembersSvc service.RoomMembersSvc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIDStr := r.PathValue("room_id")
		userID := auth.UserIDFromContext(r.Context())

		if userID == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			http.Error(w, "invalid room_id", http.StatusBadRequest)
			return
		}

		isMember, err := roomMembersSvc.IsMember(
			r.Context(),
			&domain.RoomMembers{RoomID: roomID, UserID: userID},
		)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if !isMember {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("websocket upgrade failed", "error", err)
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
