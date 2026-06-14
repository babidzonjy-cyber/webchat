package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"web-chat/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

type userMemory struct {
	users  map[int]*domain.User
	nextID int
	rwmtx  sync.RWMutex
}

func NewUserMemory() *userMemory {
	return &userMemory{
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

func (u *userMemory) Create(ctx context.Context, user *domain.User) error {
	u.rwmtx.Lock()
	defer u.rwmtx.Unlock()

	user.ID = u.nextID
	u.nextID++
	user.CreatedAt = time.Now()
	u.users[user.ID] = user

	return nil
}

func (u *userMemory) GetByID(ctx context.Context, id int) (*domain.User, error) {
	u.rwmtx.RLock()
	defer u.rwmtx.RUnlock()

	if val, exists := u.users[id]; exists {
		return val, nil
	}
	return nil, errors.New("there is no user with that id")
}

func (u *userMemory) Update(ctx context.Context, user *domain.User) error {
	u.rwmtx.Lock()
	defer u.rwmtx.Unlock()

	if _, exists := u.users[user.ID]; !exists {
		return fmt.Errorf("user %d not found", user.ID)
	}

	u.users[user.ID] = user
	return nil
}

func (u *userMemory) Delete(ctx context.Context, id int) error {
	u.rwmtx.Lock()
	defer u.rwmtx.Unlock()

	if _, exists := u.users[id]; exists {
		delete(u.users, id)
		return nil
	}

	return errors.New("there is no user with that id")
}
