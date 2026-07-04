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

type mockUserHandler struct {
	createFunc  func(ctx context.Context, user *domain.User) error
	getByIDFunc func(ctx context.Context, id int) (*domain.User, error)
	updateFunc  func(ctx context.Context, user *domain.User) error
	deleteFunc  func(ctx context.Context, id int) error

	authenticateFunc func(ctx context.Context, email, password string) (*domain.User, error)
}

func (m *mockUserHandler) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil
}

func (m *mockUserHandler) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserHandler) Update(ctx context.Context, user *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return nil
}

func (m *mockUserHandler) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockUserHandler) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(ctx, email, password)
	}
	return nil, nil
}

func TestUserHandler_Create_Success(t *testing.T) {
	mock := &mockUserHandler{}
	h := NewUserHandler(mock)

	body := `
	{"full_name":"TestUser", "email":"test@teste.com", "password":"1234567"}
	`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestUserHandler_Create_BadJSON(t *testing.T) {
	mock := &mockUserHandler{}
	h := NewUserHandler(mock)

	body := `
	{gdgsig broke
	`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_GetByID_Success(t *testing.T) {
	mock := &mockUserHandler{
		getByIDFunc: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{
				ID:       1,
				FullName: "test",
				Email:    "test@test.com",
			}, nil
		},
	}
	h := NewUserHandler(mock)

	body := `
	{"id":1,"full_name":"test","email":"test@test.com","password":"12345678"}
	`

	req := httptest.NewRequest("GET", "/users/1", strings.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUserHandler_Update_Success(t *testing.T) {
	mock := &mockUserHandler{}
	h := NewUserHandler(mock)

	body := `
	{"id":1,"full_name":"test","email":"test@test.com","password":"12345678"}
	`

	req := httptest.NewRequest("PUT", "/users/1", strings.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), auth.UserIDKey, 1)
	req = req.WithContext(ctx)

	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUserHandler_Update_BadJSON(t *testing.T) {
	mock := &mockUserHandler{}
	h := NewUserHandler(mock)

	body := `
	{"id":1,"ful eegoewrg goijdg
	`

	req := httptest.NewRequest("PUT", "/users/1", strings.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Delete_Success(t *testing.T) {
	mock := &mockUserHandler{}
	h := NewUserHandler(mock)

	body := `
	{"id":1,"full_name":"test","email":"test@test.com","password":"12345678"}
	`

	req := httptest.NewRequest("DELETE", "/users/1", strings.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), auth.UserIDKey, 1)
	req = req.WithContext(ctx)

	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}
