package service

import (
	"context"
	"fmt"
	"web-chat/internal/domain"
	"web-chat/internal/repository"
)

type RoomMembersSvc interface {
	IsMember(ctx context.Context, members *domain.RoomMembers) (bool, error)
	Add(ctx context.Context, members *domain.RoomMembers) error
	Remove(ctx context.Context, members *domain.RoomMembers) error
	GetByRoom(ctx context.Context, roomID int) ([]*domain.RoomMembers, error)
}

type roomMembersMemory struct {
	repo repository.RoomMembersRepo
}

func NewRoomMembersMemory(repo repository.RoomMembersRepo) *roomMembersMemory {
	return &roomMembersMemory{
		repo: repo,
	}
}

func (rm *roomMembersMemory) Add(ctx context.Context, members *domain.RoomMembers) error {
	if err := rm.repo.Add(ctx, members); err != nil {
		return fmt.Errorf("roomMembers.Add %v, error: %w", members, err)
	}

	return nil
}

func (rm *roomMembersMemory) Remove(ctx context.Context, members *domain.RoomMembers) error {
	if err := rm.repo.Remove(ctx, members); err != nil {
		return fmt.Errorf("roomMembers.Remove %v, error: %w", members, err)
	}

	return nil
}

func (rm *roomMembersMemory) GetByRoom(ctx context.Context, roomID int) ([]*domain.RoomMembers, error) {
	roomMembers, err := rm.repo.GetByRoom(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("roomMembers.GetByRoom, roomID %d, error %w", roomID, err)
	}

	return roomMembers, nil
}

func (rm *roomMembersMemory) IsMember(ctx context.Context, members *domain.RoomMembers) (bool, error) {
	isMember, err := rm.repo.IsMember(ctx, members)
	if err != nil {
		return false, fmt.Errorf("roomMembers.IsMember %v, error %w", members, err)
	}

	return isMember, nil
}
