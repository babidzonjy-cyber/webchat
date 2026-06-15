package repository

import (
	"context"
	"web-chat/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

type userPG struct {
	pool *pgxpool.Pool
}

func NewUserPG(pool *pgxpool.Pool) *userPG {
	return &userPG{pool: pool}
}

func (u *userPG) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO webchat.users(full_name, email, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return u.pool.QueryRow(ctx, query, user.FullName, user.Email, user.Password, user.CreatedAt).Scan(&user.ID, &user.CreatedAt)
}

func (u *userPG) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, full_name, email, password, created_at
	FROM webchat.users WHERE id = $1`

	user := &domain.User{}
	if err := u.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.FullName, &user.Email, &user.Password, &user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userPG) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE webchat.users
	SET full_name = $1, email = $2, password = $3
	WHERE id = $4`

	_, err := u.pool.Exec(ctx, query, user.FullName, user.Email, user.Password, user.ID)
	return err
}

func (u *userPG) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM webchat.users WHERE id = $1`
	_, err := u.pool.Exec(ctx, query, id)

	return err
}
