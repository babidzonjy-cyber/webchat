package repository

import (
	"context"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomMembersRepo interface {
	IsMember(ctx context.Context, members *domain.RoomMembers) (bool, error)
	Add(ctx context.Context, members *domain.RoomMembers) error
	Remove(ctx context.Context, members *domain.RoomMembers) error
	GetByRoom(ctx context.Context, roomID int) ([]*domain.RoomMembers, error)
}

type roomMembersPG struct {
	pool *pgxpool.Pool
}

func NewRoomMembersPG(pool *pgxpool.Pool) *roomMembersPG {
	return &roomMembersPG{pool: pool}
}

func (rm *roomMembersPG) Add(ctx context.Context, members *domain.RoomMembers) error {
	query := `
		INSERT INTO webchat.room_members(room_id, user_id)
		VALUES($1, $2)
		RETURNING created_at
	`

	err := rm.pool.QueryRow(
		ctx,
		query,
		members.RoomID,
		members.UserID,
	).Scan(
		&members.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (rm *roomMembersPG) Remove(ctx context.Context, members *domain.RoomMembers) error {
	query := `
		DELETE FROM webchat.room_members
		WHERE room_id = $1 AND user_id = $2
	`

	tag, err := rm.pool.Exec(
		ctx,
		query,
		members.RoomID,
		members.UserID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (rm *roomMembersPG) GetByRoom(ctx context.Context, roomID int) ([]*domain.RoomMembers, error) {
	query := `
		SELECT room_id, user_id, created_at FROM webchat.room_members
		WHERE room_id = $1
	`

	rows, err := rm.pool.Query(
		ctx,
		query,
		roomID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.RoomMembers

	for rows.Next() {
		member := &domain.RoomMembers{}
		if err := rows.Scan(
			&member.RoomID,
			&member.UserID,
			&member.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, rows.Err()
}

func (rm *roomMembersPG) IsMember(ctx context.Context, members *domain.RoomMembers) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM webchat.room_members
			WHERE room_id = $1 AND user_id = $2
		)
	`

	var isMember bool

	err := rm.pool.QueryRow(
		ctx,
		query,
		members.RoomID,
		members.UserID,
	).Scan(
		&isMember,
	)
	if err != nil {
		return false, err
	}

	return isMember, nil
}
