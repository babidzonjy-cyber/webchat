package service

import (
	"context"
	"fmt"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"
	"web-chat/internal/repository"
)

type MessageService interface {
	Create(ctx context.Context, msg *domain.Message) error
	GetByID(ctx context.Context, id int) (*domain.Message, error)
	GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error)
	Delete(ctx context.Context, room_id, userID int) error      // только автор удаляет свое сообщение
	DeleteByRoom(ctx context.Context, roomID, userID int) error // только создатель комнаты удаляет все сообщения в группе
}

type messageMemory struct {
	repo     repository.MessageRepository
	roomRepo repository.RoomRepository
}

func NewMessageMemory(repo repository.MessageRepository, roomRepo repository.RoomRepository) *messageMemory {
	return &messageMemory{
		repo:     repo,
		roomRepo: roomRepo,
	}
}

func (m *messageMemory) Create(ctx context.Context, msg *domain.Message) error {
	if err := m.repo.Create(ctx, msg); err != nil {
		return fmt.Errorf("service.Create msg: %d, error: %w", msg.ID, err)
	}

	return nil
}

func (m *messageMemory) GetByID(ctx context.Context, id int) (*domain.Message, error) {
	msg, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID msg: %d, error: %w", id, err)
	}

	return msg, nil
}

func (m *messageMemory) GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error) {
	msgs, err := m.repo.GetByRoomID(ctx, roomID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service.GetByRoomID msgs %v, roomID %d, error: %w", msgs, roomID, err)
	}

	return msgs, nil
}

func (m *messageMemory) Delete(ctx context.Context, msgID, userID int) error {
	if err := m.repo.Delete(ctx, msgID, userID); err != nil {
		return fmt.Errorf("service.Delete msg %d, userID %d, error: %w", msgID, userID, err)
	}

	return nil
}

func (m *messageMemory) DeleteByRoom(ctx context.Context, roomID, userID int) error {
	room, err := m.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return apperrors.ErrNotFound
	}
	if room.CreatedBy != userID {
		return apperrors.ErrForbidden
	}

	if err := m.repo.DeleteByRoom(ctx, roomID); err != nil {
		return fmt.Errorf("service.DeleteByRoom msgs, userID %d, roomID %d, error: %w", userID, roomID, err)
	}

	return nil
}
