package websocket

import "time"

type IncomingMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type OutgoingMessage struct {
	Type      string    `json:"type"`
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
