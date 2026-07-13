package service

import (
	"context"
	"web-chat/internal/domain"
)

type mockRoomMembersService struct {
	isMemberFunc  func(ctx context.Context, members *domain.RoomMembers) (bool, error)
	addFunc       func(ctx context.Context, members *domain.RoomMembers) error
	removeFunc    func(ctx context.Context, members *domain.RoomMembers) error
	getByRoomFunc func(ctx context.Context, roomID int) ([]*domain.RoomMembers, error)
}

func (rm *mockRoomMembersService) IsMember(ctx context.Context, members *domain.RoomMembers) (bool, error) {
	if rm.isMemberFunc != nil {
		return rm.isMemberFunc(ctx, members)
	}
	return true, nil
}

func (rm *mockRoomMembersService) Add(ctx context.Context, members *domain.RoomMembers) error {
	if rm.addFunc != nil {
		return rm.addFunc(ctx, members)
	}
	return nil
}

func (rm *mockRoomMembersService) Remove(ctx context.Context, members *domain.RoomMembers) error {
	if rm.removeFunc != nil {
		return rm.removeFunc(ctx, members)
	}
	return nil
}

func (rm *mockRoomMembersService) GetByRoom(ctx context.Context, roomID int) ([]*domain.RoomMembers, error) {
	if rm.getByRoomFunc != nil {
		return rm.getByRoomFunc(ctx, roomID)
	}
	return nil, nil
}
