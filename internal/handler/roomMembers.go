package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-chat/internal/domain"
	"web-chat/internal/dto"
	"web-chat/internal/service"
)

type RoomMembersHandler struct {
	svc service.RoomMembersSvc
}

type IsMember struct {
	IsMember bool `json:"is_member"`
}

func NewRoomMembersHandler(svc service.RoomMembersSvc) *RoomMembersHandler {
	return &RoomMembersHandler{
		svc: svc,
	}
}

func (rm *RoomMembersHandler) Add(w http.ResponseWriter, r *http.Request) {
	var members domain.RoomMembers

	roomID, err := strconv.Atoi(r.PathValue("room_id"))
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}
	members.RoomID = roomID

	if err := json.NewDecoder(r.Body).Decode(&members); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rm.svc.Add(r.Context(), &members); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(roomMembersToResponse(&members))
}

func (rm *RoomMembersHandler) Remove(w http.ResponseWriter, r *http.Request) {
	roomIDstr := r.PathValue("room_id")
	userIDStr := r.PathValue("user_id")

	roomID, err := strconv.Atoi(roomIDstr)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	member := &domain.RoomMembers{
		RoomID: roomID,
		UserID: userID,
	}

	if err := rm.svc.Remove(r.Context(), member); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (rm *RoomMembersHandler) List(w http.ResponseWriter, r *http.Request) {
	roomIDstr := r.PathValue("room_id")

	roomID, err := strconv.Atoi(roomIDstr)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	members, err := rm.svc.GetByRoom(r.Context(), roomID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(membersToResponse(members))
}

func (rm *RoomMembersHandler) IsMember(w http.ResponseWriter, r *http.Request) {
	roomIDstr := r.PathValue("room_id")
	userIDstr := r.PathValue("user_id")

	roomID, err := strconv.Atoi(roomIDstr)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	member := &domain.RoomMembers{
		RoomID: roomID,
		UserID: userID,
	}

	ok, err := rm.svc.IsMember(r.Context(), member)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(IsMember{
		IsMember: ok,
	})

	w.WriteHeader(http.StatusOK)
}

func roomMembersToResponse(rm *domain.RoomMembers) dto.RoomMembers {
	return dto.RoomMembers{
		RoomID:    rm.RoomID,
		UserID:    rm.UserID,
		CreatedAt: rm.CreatedAt,
	}
}

func membersToResponse(members []*domain.RoomMembers) []dto.RoomMembers {
	out := make([]dto.RoomMembers, 0, len(members))
	for _, v := range members {
		out = append(out, roomMembersToResponse(v))
	}
	return out
}
