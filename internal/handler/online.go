package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-chat/internal/service"
)

type OnlineHandler struct {
	svc service.OnlineService
}

func NewOnlineHandler(svc service.OnlineService) *OnlineHandler {
	return &OnlineHandler{
		svc: svc,
	}
}

type OnlineCountResponse struct {
	Count int `json:"count"`
}

type CheckOnlineResponse struct {
	Online bool `json:"online"`
}

type UsersOnlineResponse struct {
	Users []int `json:"users"`
}

func (o *OnlineHandler) OnlineCount(w http.ResponseWriter, r *http.Request) {
	strRoomID := r.PathValue("id")

	roomID, err := strconv.Atoi(strRoomID)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	count, err := o.svc.OnlineCount(roomID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&OnlineCountResponse{
		Count: count,
	})
}

func (o *OnlineHandler) CheckOnline(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.PathValue("room_id")
	userIDStr := r.PathValue("user_id")

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	ok, err := o.svc.CheckOnline(userID, roomID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&CheckOnlineResponse{
		Online: ok,
	})
}

func (o *OnlineHandler) OnlineUsers(w http.ResponseWriter, r *http.Request) {
	strRoomID := r.PathValue("id")

	roomID, err := strconv.Atoi(strRoomID)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	users, err := o.svc.OnlineUsers(roomID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&UsersOnlineResponse{
		Users: users,
	})
}
