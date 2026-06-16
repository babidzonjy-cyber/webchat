package service

import (
	"context"
	"web-chat/internal/apperrors"
	"web-chat/internal/domain"
	"web-chat/internal/repository"
)

type MessageService interface {
	Create(ctx context.Context, msg *domain.Message) error
	GetByID(ctx context.Context, id int) (*domain.Message, error)
	GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error)
	Delete(ctx context.Context, id, userID int) error           // только автор удаляет свое сообщение
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
		return err
	}

	return nil
}

func (m *messageMemory) GetByID(ctx context.Context, id int) (*domain.Message, error) {
	msg, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *messageMemory) GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error) {
	msgs, err := m.repo.GetByRoomID(ctx, roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	return msgs, err
}

func (m *messageMemory) Delete(ctx context.Context, msgID, userID int) error {
	if err := m.repo.Delete(ctx, msgID, userID); err != nil {
		return err
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
	return m.repo.DeleteByRoom(ctx, roomID)
}
