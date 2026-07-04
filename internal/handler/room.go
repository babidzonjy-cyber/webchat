package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-chat/internal/auth"
	"web-chat/internal/domain"
	"web-chat/internal/dto"
	"web-chat/internal/service"
)

type RoomHandler struct {
	svc service.RoomService
}

func NewRoomHandler(svc service.RoomService) *RoomHandler {
	return &RoomHandler{
		svc: svc,
	}
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var room domain.Room

	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	room.CreatedBy = auth.UserIDFromContext(r.Context())
	if err := h.svc.Create(r.Context(), &room); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(roomToResponse(&room))
}

func (h *RoomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	room, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(roomToResponse(room))
}

func (h *RoomHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.svc.GetAll(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(roomsToResponse(rooms))
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	var room domain.Room
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	room.ID = id
	existing, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}

	if existing.CreatedBy != auth.UserIDFromContext(r.Context()) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.svc.Update(r.Context(), &room); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(roomToResponse(&room))
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	existing, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if existing.CreatedBy != auth.UserIDFromContext(r.Context()) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func roomToResponse(r *domain.Room) dto.RoomDTO {
	return dto.RoomDTO{
		ID:        r.ID,
		Name:      r.Name,
		CreatedBy: r.CreatedBy,
		CreatedAt: r.CreatedAt,
	}
}

func roomsToResponse(rooms []*domain.Room) []dto.RoomDTO {
	out := make([]dto.RoomDTO, 0, len(rooms))
	for _, r := range rooms {
		out = append(out, roomToResponse(r))
	}
	return out
}
