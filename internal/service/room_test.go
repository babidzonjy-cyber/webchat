package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"
)

type mockRoomService struct {
	createFunc  func(ctx context.Context, room *domain.Room) error
	getByIDFunc func(ctx context.Context, id int) (*domain.Room, error)
	getAllFunc  func(ctx context.Context) ([]*domain.Room, error)
	updateFunc  func(ctx context.Context, room *domain.Room) error
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *mockRoomService) Create(ctx context.Context, room *domain.Room) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, room)
	}
	return nil
}

func (m *mockRoomService) GetByID(ctx context.Context, id int) (*domain.Room, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRoomService) GetAll(ctx context.Context) ([]*domain.Room, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockRoomService) Update(ctx context.Context, room *domain.Room) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, room)
	}
	return nil
}

func (m *mockRoomService) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestCreateRoom_Success(t *testing.T) {
	mock := &mockRoomService{}
	svc := NewRoomMemory(mock)

	room := &domain.Room{
		Name: "name",
	}

	err := svc.Create(context.Background(), room)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestGetByIDRoom_Success(t *testing.T) {
	mock := &mockRoomService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Room, error) {
			return &domain.Room{
				ID:   1,
				Name: "name",
			}, nil
		},
	}

	svc := NewRoomMemory(mock)

	room, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
	if room == nil {
		t.Errorf("user shouldn't be nil")
	}

	fmt.Printf("user: %v\n", room)
}

func TestGetAllRoom_Success(t *testing.T) {
	mock := &mockRoomService{
		getAllFunc: func(ctx context.Context) ([]*domain.Room, error) {
			return []*domain.Room{
				{ID: 1, Name: "room1"},
				{ID: 2, Name: "room2"},
				{ID: 3, Name: "room3"},
				{ID: 4, Name: "room4"},
			}, nil
		},
	}
	svc := NewRoomMemory(mock)

	rooms, err := svc.GetAll(context.Background())
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}

	if len(rooms) != 4 {
		t.Errorf("expected 4 rooms, got %d", len(rooms))
	}
}

func TestUpdateRoom_Success(t *testing.T) {
	mock := &mockRoomService{}

	svc := NewRoomMemory(mock)

	room := &domain.Room{
		Name: "name",
	}

	err := svc.Update(context.Background(), room)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestDeleteRoom_Success(t *testing.T) {
	mock := &mockRoomService{}

	svc := NewRoomMemory(mock)

	err := svc.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	mock := &mockRoomService{
		deleteFunc: func(ctx context.Context, id int) error {
			return apperrors.ErrNotFound
		},
	}

	svc := NewRoomMemory(mock)

	err := svc.Delete(context.Background(), 999)
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}
