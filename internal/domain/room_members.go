package domain

import "time"

type RoomMembers struct {
	RoomID    int       `json:"room_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
