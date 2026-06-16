package service

import (
	"context"
	"errors"
	"fmt"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"
	"web-chat/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

type userMemory struct {
	repo repository.UserRepository
}

func NewUserMemory(repo repository.UserRepository) *userMemory {
	return &userMemory{
		repo: repo,
	}
}

func (u *userMemory) Create(ctx context.Context, user *domain.User) error {
	if err := u.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userMemory) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *userMemory) Update(ctx context.Context, user *domain.User) error {
	if err := u.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userMemory) Delete(ctx context.Context, id int) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("service.Delete user %d: %w", id, err)
	}

	return nil
}
