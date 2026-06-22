package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"
	"web-chat/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error

	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
}

type userMemory struct {
	repo repository.UserRepository
}

func NewUserMemory(repo repository.UserRepository) *userMemory {
	return &userMemory{
		repo: repo,
	}
}

func validateUser(user *domain.User) error {
	if user.FullName == "" {
		return fmt.Errorf("%w: full_name is required", apperrors.ErrValidation)
	}

	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("%w: invalid email", apperrors.ErrValidation)
	}

	if len(user.Password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", apperrors.ErrValidation)
	}

	return nil
}

func (u *userMemory) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("service.Authenticate: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, apperrors.ErrNotFound
	}

	return user, nil
}

func (u *userMemory) Create(ctx context.Context, user *domain.User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("service.Create hash password: %w", err)
	}

	user.Password = string(hashedPassword)

	if err := u.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("service.Create user: %d, error: %w", user.ID, err)
	}

	return nil
}

func (u *userMemory) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetById %d: %w", id, err)
	}

	return user, nil
}

func (u *userMemory) Update(ctx context.Context, user *domain.User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	if err := u.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("service.Update user: %d, error: %w", user.ID, err)
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
