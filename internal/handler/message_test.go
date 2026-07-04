package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"web-chat/internal/domain"
)

type mockMessageHandler struct {
	createFunc       func(ctx context.Context, msg *domain.Message) error
	getByIDFunc      func(ctx context.Context, id int) (*domain.Message, error)
	getByRoomIDFunc  func(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error)
	deleteFunc       func(ctx context.Context, msgID, userID int) error
	deleteByRoomFunc func(ctx context.Context, roomID, userID int) error
}

func (m *mockMessageHandler) Create(ctx context.Context, msg *domain.Message) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, msg)
	}
	return nil
}

func (m *mockMessageHandler) GetByID(ctx context.Context, id int) (*domain.Message, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockMessageHandler) GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error) {
	if m.getByRoomIDFunc != nil {
		return m.getByRoomIDFunc(ctx, roomID, limit, offset)
	}
	return nil, nil
}

func (m *mockMessageHandler) Delete(ctx context.Context, room_id, userID int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, room_id, userID)
	}
	return nil
}

func (m *mockMessageHandler) DeleteByRoom(ctx context.Context, roomID, userID int) error {
	if m.deleteByRoomFunc != nil {
		return m.deleteByRoomFunc(ctx, roomID, userID)
	}
	return nil
}

func TestMessageHandler_Create_Success(t *testing.T) {
	mock := &mockMessageHandler{}
	h := NewMessageHandler(mock)

	body := `
	{"text":"testtext","room_id":1,"user_id":1}
	`

	req := httptest.NewRequest("POST", "/rooms/{room_id}/messages", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestMessageHandler_Creat_BadJSON(t *testing.T) {
	mock := &mockMessageHandler{}
	h := NewMessageHandler(mock)

	body := `
	{"text":"test dgsger reger
	`

	req := httptest.NewRequest("POST", "/rooms/{room_id}/messages", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestMessageHandler_GetByRoomID_Success(t *testing.T) {
	mock := &mockMessageHandler{}
	h := NewMessageHandler(mock)

	req := httptest.NewRequest("GET", "/rooms/1/messages?limit=10&offset=0", nil)
	req.SetPathValue("room_id", "1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.GetByRoomID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMessageHandler_GetByID_Success(t *testing.T) {
	mock := &mockMessageHandler{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Message, error) {
			return &domain.Message{
				ID:     1,
				Text:   "text",
				RoomID: 1,
				UserID: 1,
			}, nil
		},
	}
	h := NewMessageHandler(mock)

	body := `
	{"text":"testtext","room_id":1,"user_id":1}
	`

	req := httptest.NewRequest("GET", "/messages/1", strings.NewReader(body))
	req.SetPathValue("id", "1")

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMessageHandler_Delete_Success(t *testing.T) {
	mock := &mockMessageHandler{}
	h := NewMessageHandler(mock)

	req := httptest.NewRequest("DELETE", "/messages/1?user_id=1", nil)
	req.SetPathValue("id", "1")

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestMessageHandler_DeleteByRoom_Success(t *testing.T) {
	mock := &mockMessageHandler{}
	h := NewMessageHandler(mock)

	req := httptest.NewRequest("DELETE", "/rooms/1/messages?user_id=1", nil)
	req.SetPathValue("room_id", "1")

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.DeleteByRoom(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}
