package handler

import (
	"encoding/json"
	"net/http"
	"web-chat/internal/auth"
	"web-chat/internal/domain"
	"web-chat/internal/service"
)

type AuthHandler struct {
	svc service.UserService
}

func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.Create(r.Context(), &user); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(userToResponse(&user))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	creds := loginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Authenticate(r.Context(), creds.Email, creds.Password)
	if err != nil {
		writeAppError(w, err)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "generate token error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	}); err != nil {
		http.Error(w, "cannot encode data", http.StatusInternalServerError)
		return
	}
}
