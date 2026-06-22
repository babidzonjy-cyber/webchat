package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"web-chat/internal/apperrors"
)

func writeAppError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, apperrors.ErrForbidden):
		http.Error(w, "forbidden", http.StatusForbidden)
	case errors.Is(err, apperrors.ErrValidation):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrConflict):
		http.Error(w, "user with this email already exists", http.StatusConflict)
	default:
		slog.Error("internal error", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
