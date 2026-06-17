package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"
)

type mockUserService struct {
	createFunc  func(ctx context.Context, user *domain.User) error
	getByIDFunc func(ctx context.Context, id int) (*domain.User, error)
	updateFunc  func(ctx context.Context, user *domain.User) error
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *mockUserService) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}

func (m *mockUserService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserService) Update(ctx context.Context, user *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return nil
}

func (m *mockUserService) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestCreateUser_Success(t *testing.T) {
	mock := &mockUserService{}
	svc := NewUserMemory(mock)
	user := &domain.User{
		FullName: "name",
		Email:    "2@2sdsdsd.ru",
	}

	err := svc.Create(context.Background(), user)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestGetByIDUser_Success(t *testing.T) {
	mock := &mockUserService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				FullName: "name",
				Email:    "tiger@gmail.com",
			}, nil
		},
	}

	svc := NewUserMemory(mock)

	user, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
	if user == nil {
		t.Errorf("user shouldn't be nil")
	}

	fmt.Printf("user: %v\n", user)
}

func TestUpdateUser_Success(t *testing.T) {
	mock := &mockUserService{}

	svc := NewUserMemory(mock)

	user := &domain.User{
		FullName: "name",
		Email:    "2@2sdsdsd.ru",
	}

	err := svc.Update(context.Background(), user)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mock := &mockUserService{}

	svc := NewUserMemory(mock)

	err := svc.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mock := &mockUserService{
		deleteFunc: func(ctx context.Context, id int) error {
			return apperrors.ErrNotFound
		},
	}

	svc := NewUserMemory(mock)

	err := svc.Delete(context.Background(), 999)
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
