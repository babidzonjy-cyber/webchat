package service

import (
	"context"
	"testing"
	"web-chat/internal/domain"
)

type mockMessageService struct {
	createFunc       func(ctx context.Context, msg *domain.Message) error
	getByIDFunc      func(ctx context.Context, id int) (*domain.Message, error)
	getByRoomIDFunc  func(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error)
	deleteFunc       func(ctx context.Context, msgID, userID int) error
	deleteByRoomFunc func(ctx context.Context, roomID int) error
}

func (m *mockMessageService) Create(ctx context.Context, msg *domain.Message) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, msg)
	}
	return nil
}

func (m *mockMessageService) GetByID(ctx context.Context, id int) (*domain.Message, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockMessageService) GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error) {
	if m.getByRoomIDFunc != nil {
		return m.getByRoomIDFunc(ctx, roomID, limit, offset)
	}
	return nil, nil
}

func (m *mockMessageService) Delete(ctx context.Context, room_id, userID int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, room_id, userID)
	}
	return nil
}

func (m *mockMessageService) DeleteByRoom(ctx context.Context, roomID int) error {
	if m.deleteByRoomFunc != nil {
		return m.deleteByRoomFunc(ctx, roomID)
	}
	return nil
}

func TestCreateMessage_Success(t *testing.T) {
	mockMsg := &mockMessageService{}
	mockRoom := &mockRoomService{}
	svc := NewMessageMemory(mockMsg, mockRoom)

	msg := &domain.Message{
		Text:   "text",
		RoomID: 1,
		UserID: 1,
	}

	err := svc.Create(context.Background(), msg)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestGetByIDMessage_Success(t *testing.T) {
	mockMsg := &mockMessageService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Message, error) {
			return &domain.Message{
				ID:     1,
				RoomID: 1,
				UserID: 1,
				Text:   "text",
			}, nil
		},
	}

	mockRoom := &mockRoomService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Room, error) {
			return &domain.Room{
				ID:        1,
				CreatedBy: 1,
			}, nil
		},
	}
	svc := NewMessageMemory(mockMsg, mockRoom)

	msg, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}

	if msg.ID != 1 {
		t.Errorf("got %d, want %d", msg.ID, 1)
	}

	if msg.RoomID != 1 {
		t.Errorf("got %d, want %d", msg.RoomID, 1)
	}

	if msg.UserID != 1 {
		t.Errorf("got %d, want %d", msg.UserID, 1)
	}
}

func TestGetByRoomIDMessage_Success(t *testing.T) {
	mockMsg := &mockMessageService{}
	mockRoom := &mockRoomService{}
	svc := NewMessageMemory(mockMsg, mockRoom)

	_, err := svc.GetByRoomID(context.Background(), 1, 10, 0)
	if err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestDeleteMessage_Success(t *testing.T) {
	mockMsg := &mockMessageService{}
	mockRoom := &mockRoomService{}
	svc := NewMessageMemory(mockMsg, mockRoom)

	if err := svc.Delete(context.Background(), 1, 1); err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}

func TestDeleteByRoomMessage_Success(t *testing.T) {
	mockMsg := &mockMessageService{}
	mockRoom := &mockRoomService{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Room, error) {
			return &domain.Room{
				ID:        1,
				CreatedBy: 1,
			}, nil
		},
	}
	svc := NewMessageMemory(mockMsg, mockRoom)

	if err := svc.DeleteByRoom(context.Background(), 1, 1); err != nil {
		t.Errorf("не должно быть ошибки, а получили: %v", err)
	}
}
