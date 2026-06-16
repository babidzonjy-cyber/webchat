package repository

import (
	"context"
	"fmt"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	GetByID(ctx context.Context, id int) (*domain.Message, error)
	GetByRoomID(ctx context.Context, roomID, limit, offset int) ([]*domain.Message, error)
	Delete(ctx context.Context, msgID, userID int) error
	DeleteByRoom(ctx context.Context, roomID int) error
}

type messagePG struct {
	pool *pgxpool.Pool
}

func NewMessagePG(pool *pgxpool.Pool) *messagePG {
	return &messagePG{pool: pool}
}

func (m *messagePG) Create(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO webchat.messages(text, room_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return m.pool.QueryRow(ctx, query, message.Text, message.RoomID, message.UserID, message.CreatedAt).Scan(&message.ID, &message.CreatedAt)
}

func (m *messagePG) GetByID(ctx context.Context, id int) (*domain.Message, error) {
	query := `SELECT id, text, user_id, room_id, created_at
	FROM webchat.messages WHERE id = $1`

	message := &domain.Message{}
	if err := m.pool.QueryRow(ctx, query, id).Scan(
		&message.ID, &message.Text, &message.UserID, &message.RoomID, &message.CreatedAt); err != nil {
		return nil, err
	}

	return message, nil
}

func (m *messagePG) GetByRoomID(ctx context.Context, roomID, limit, offset int) ([]*domain.Message, error) {
	query := `SELECT id, text, room_id, user_id, created_at
	FROM webchat.messages
	WHERE room_id = $1
	ORDER BY created_at ASC
	LIMIT $2 OFFSET $3`

	rows, err := m.pool.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*domain.Message

	for rows.Next() {
		msg := &domain.Message{}
		if err := rows.Scan(&msg.ID, &msg.Text, &msg.RoomID, &msg.UserID, &msg.CreatedAt); err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, rows.Err()
}

func (m *messagePG) Delete(ctx context.Context, msgID, userID int) error {
	query := `DELETE FROM webchat.messages WHERE id = $1 and user_id = $2`
	tag, err := m.pool.Exec(ctx, query, msgID, userID)
	if err != nil {
		return fmt.Errorf("repo.Delete, msgId %d, userId %d: %w", msgID, userID, err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (m *messagePG) DeleteByRoom(ctx context.Context, roomID int) error {
	query := `DELETE FROM webchat.messages WHERE room_id = $1`
	_, err := m.pool.Exec(ctx, query, roomID)

	return err
}
