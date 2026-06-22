package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"web-chat/internal/auth"
	"web-chat/internal/domain"
)

type mockRoomHandler struct {
	createFunc  func(ctx context.Context, room *domain.Room) error
	getByIDFunc func(ctx context.Context, id int) (*domain.Room, error)
	getAllFunc  func(ctx context.Context) ([]*domain.Room, error)
	updateFunc  func(ctx context.Context, room *domain.Room) error
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *mockRoomHandler) Create(ctx context.Context, room *domain.Room) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, room)
	}
	return nil
}

func (m *mockRoomHandler) GetByID(ctx context.Context, id int) (*domain.Room, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRoomHandler) GetAll(ctx context.Context) ([]*domain.Room, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockRoomHandler) Update(ctx context.Context, room *domain.Room) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, room)
	}
	return nil
}

func (m *mockRoomHandler) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestRoomHandler_Create_Success(t *testing.T) {
	mock := &mockRoomHandler{}
	h := NewRoomHandler(mock)

	body := `
	{"name":"TestRoom"}
	`

	req := httptest.NewRequest("POST", "/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestRoomHandler_Create_BadJSON(t *testing.T) {
	mock := &mockRoomHandler{}
	h := NewRoomHandler(mock)

	body := `
	{"name  eege
	`

	req := httptest.NewRequest("POST", "/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRoomHandler_GetByID_Success(t *testing.T) {
	mock := &mockRoomHandler{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Room, error) {
			return &domain.Room{
				Name: "TestRoom",
			}, nil
		},
	}
	h := NewRoomHandler(mock)

	body := `
	{"id":1,"name":"TestRoom1"}
	`

	req := httptest.NewRequest("GET", "/rooms/1", strings.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRoomHandler_GetAll_Success(t *testing.T) {
	mock := &mockRoomHandler{}
	h := NewRoomHandler(mock)

	body := `
	[
		{"id":1,"name":"room1"},
		{"id":2,"name":"room2"},
		{"id":3,"name":"room3"}
	]
	`

	req := httptest.NewRequest("GET", "/rooms", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRoomHandler_Delete_Success(t *testing.T) {
	mock := &mockRoomHandler{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Room, error) {
			return &domain.Room{
				ID:        1,
				Name:      "room1",
				CreatedBy: 1,
			}, nil
		},
	}
	h := NewRoomHandler(mock)

	body := `
	{"id":1,"name":"room1"}
	`

	req := httptest.NewRequest("DELETE", "/rooms/1", strings.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), auth.UserIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}
