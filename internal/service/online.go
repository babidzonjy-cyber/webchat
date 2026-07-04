package service

import (
	"fmt"
	"web-chat/internal/repository"
)

type OnlineService interface {
	OnlineCount(roomID int) (int, error)
	CheckOnline(userID, roomID int) (bool, error)
	OnlineUsers(roomID int) ([]int, error)
}

type onlineMemory struct {
	svc repository.OnlineRepository
}

func NewOnlineService(svc repository.OnlineRepository) *onlineMemory {
	return &onlineMemory{
		svc: svc,
	}
}

func (o *onlineMemory) OnlineCount(roomID int) (int, error) {
	res, err := o.svc.GetOnlineCount(roomID)
	if err != nil {
		return -1, fmt.Errorf("cannot count users in room %d, error: %w", roomID, err)
	}
	return res, nil
}

func (o *onlineMemory) CheckOnline(userID, roomID int) (bool, error) {
	ok, err := o.svc.IsOnline(roomID, userID)
	if err != nil {
		return false, fmt.Errorf("cannot check online %w", err)
	}
	return ok, nil
}

func (o *onlineMemory) OnlineUsers(roomID int) ([]int, error) {
	users, err := o.svc.GetOnlineUsers(roomID)
	if err != nil {
		return nil, fmt.Errorf("cannot give online users %w", err)
	}
	return users, nil
}
