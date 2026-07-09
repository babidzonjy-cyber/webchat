package handler

import (
	"context"
	"web-chat/internal/domain"
)

type mockRoomMembersHandler struct {
	addFunc       func(ctx context.Context, members *domain.RoomMembers) error
	isMemberFunc  func(ctx context.Context, members *domain.RoomMembers) (bool, error)
	removeFunc    func(ctx context.Context, members *domain.RoomMembers) error
	getByRoomFunc func(ctx context.Context, roomID int) ([]*domain.RoomMembers, error)
}

func (m *mockRoomMembersHandler) Add(ctx context.Context, members *domain.RoomMembers) error {
	if m.addFunc != nil {
		return m.addFunc(ctx, members)
	}
	return nil
}

func (m *mockRoomMembersHandler) IsMember(ctx context.Context, members *domain.RoomMembers) (bool, error) {
	if m.isMemberFunc != nil {
		return m.isMemberFunc(ctx, members)
	}
	return true, nil
}

func (m *mockRoomMembersHandler) Remove(ctx context.Context, members *domain.RoomMembers) error {
	if m.removeFunc != nil {
		return m.removeFunc(ctx, members)
	}
	return nil
}

func (m *mockRoomMembersHandler) GetByRoom(ctx context.Context, roomID int) ([]*domain.RoomMembers, error) {
	if m.getByRoomFunc != nil {
		return m.getByRoomFunc(ctx, roomID)
	}
	return nil, nil
}
