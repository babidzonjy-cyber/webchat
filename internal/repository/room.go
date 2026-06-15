package repository

import (
	"context"
	"web-chat/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository interface {
	Create(ctx context.Context, room *domain.Room) error
	GetByID(ctx context.Context, id int) (*domain.Room, error)
	GetAll(ctx context.Context) ([]*domain.Room, error)
	Update(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, id int) error
}

type roomPG struct {
	pool *pgxpool.Pool
}

func NewRoomPG(pool *pgxpool.Pool) *roomPG {
	return &roomPG{pool: pool}
}

func (r *roomPG) Create(ctx context.Context, room *domain.Room) error {
	query := `
		INSERT INTO webchat.rooms(name, created_by,created_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.pool.QueryRow(ctx, query, room.Name, room.CreatedBy, room.CreatedAt).Scan(&room.ID, &room.CreatedAt)
}

func (r *roomPG) GetByID(ctx context.Context, id int) (*domain.Room, error) {
	query := `SELECT id, name, created_by, created_at
	FROM webchat.rooms WHERE id = $1`

	room := &domain.Room{}
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt); err != nil {
		return nil, err
	}

	return room, nil
}

func (r *roomPG) GetAll(ctx context.Context) ([]*domain.Room, error) {
	query := `SELECT id, name, created_by, created_at
	FROM webchat.rooms`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*domain.Room

	for rows.Next() {
		room := &domain.Room{}
		if err := rows.Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt); err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, rows.Err()
}

func (r *roomPG) Update(ctx context.Context, room *domain.Room) error {
	query := `UPDATE webchat.rooms SET name = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, room.Name, room.ID)

	return err
}

func (r *roomPG) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM webchat.rooms WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)

	return err
}
