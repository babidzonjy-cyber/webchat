package repository

import (
	"context"
	"errors"
	"fmt"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository interface {
	CreateWithOwner(ctx context.Context, room *domain.Room) error
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

func (r *roomPG) CreateWithOwner(ctx context.Context, room *domain.Room) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query1 := `
		INSERT INTO webchat.rooms(name, created_by)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err = tx.QueryRow(
		ctx,
		query1,
		room.Name,
		room.CreatedBy,
	).Scan(
		&room.ID,
		&room.CreatedAt,
	)
	if err != nil {
		return err
	}

	query2 := `
		INSERT INTO webchat.room_members(room_id, user_id)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(
		ctx,
		query2,
		room.ID,
		room.CreatedBy,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *roomPG) GetByID(ctx context.Context, id int) (*domain.Room, error) {
	query := `SELECT id, name, created_by, created_at
	FROM webchat.rooms WHERE id = $1`

	room := &domain.Room{}
	err := r.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&room.ID,
		&room.Name,
		&room.CreatedBy,
		&room.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
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
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.CreatedBy,
			&room.CreatedAt,
		); err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, rows.Err()
}

func (r *roomPG) Update(ctx context.Context, room *domain.Room) error {
	query := `UPDATE webchat.rooms SET name = $1 WHERE id = $2`

	_, err := r.pool.Exec(
		ctx,
		query,
		room.Name,
		room.ID,
	)

	return err
}

func (r *roomPG) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM webchat.rooms WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("repo.Delete, room %d: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
