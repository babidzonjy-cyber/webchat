package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-chat/internal/domain"
	"web-chat/internal/dto"
	"web-chat/internal/service"
)

type MessageHandler struct {
	svc service.MessageService
}

func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{
		svc: svc,
	}
}

func (m *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var message domain.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := m.svc.Create(r.Context(), &message); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(messageToResponse(&message))
}

func (m *MessageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	msg, err := m.svc.GetByID(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(messageToResponse(msg))
}

func (m *MessageHandler) GetByRoomID(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.PathValue("room_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	limit := 50
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}

	offset := 0
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}

	msgs, err := m.svc.GetByRoomID(r.Context(), roomID, limit, offset)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(messagesToResponse(msgs))
}

func (m *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	idStr := r.PathValue("id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := m.svc.Delete(r.Context(), id, userID); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (m *MessageHandler) DeleteByRoom(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	roomIDStr := r.PathValue("room_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	if err := m.svc.DeleteByRoom(r.Context(), roomID, userID); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func messageToResponse(m *domain.Message) dto.MessageDTO {
	return dto.MessageDTO{
		ID:        m.ID,
		Text:      m.Text,
		RoomID:    m.RoomID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
	}
}

func messagesToResponse(messages []*domain.Message) []dto.MessageDTO {
	out := make([]dto.MessageDTO, 0, len(messages))
	for _, m := range messages {
		out = append(out, messageToResponse(m))
	}
	return out
}
