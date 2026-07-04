package domain

import "time"

type Message struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	RoomID    int       `json:"room_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
