package service

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
	"web-chat/internal/domain"
)

type MessageService interface {
	Create(ctx context.Context, msg *domain.Message) error
	GetByID(ctx context.Context, id int) (*domain.Message, error)
	GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error)
	Delete(ctx context.Context, id, userID int) error           // только автор удаляет свое сообщение
	DeleteByRoom(ctx context.Context, roomID, userID int) error // только создатель комнаты удаляет все сообщения в группе
}

type messageMemory struct {
	message map[int]*domain.Message
	nextID  int

	rwmtx sync.RWMutex
}

func NewMessageMemory() *messageMemory {
	return &messageMemory{
		message: make(map[int]*domain.Message),
		nextID:  1,
	}
}

func (m *messageMemory) Create(ctx context.Context, msg *domain.Message) error {
	m.rwmtx.Lock()
	defer m.rwmtx.Unlock()

	msg.ID = m.nextID
	m.nextID++
	msg.CreatedAt = time.Now()
	m.message[msg.ID] = msg

	return nil
}

func (m *messageMemory) GetByID(ctx context.Context, id int) (*domain.Message, error) {
	m.rwmtx.RLock()
	defer m.rwmtx.RUnlock()

	if val, exists := m.message[id]; exists {
		return val, nil
	}

	return nil, errors.New("there is no message with that id")
}

func (m *messageMemory) GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*domain.Message, error) {
	m.rwmtx.RLock()
	defer m.rwmtx.RUnlock()

	var allMessages []*domain.Message
	for _, msg := range m.message {
		if msg.RoomID == roomID {
			allMessages = append(allMessages, msg)
		}
	}

	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].CreatedAt.Before(allMessages[j].CreatedAt)
	})

	if offset >= len(allMessages) {
		return []*domain.Message{}, nil
	}

	start := offset
	end := start + limit
	if end > len(allMessages) {
		end = len(allMessages)
	}

	return allMessages[start:end], nil
}

func (m *messageMemory) Delete(ctx context.Context, msgID, userID int) error {
	m.rwmtx.Lock()
	defer m.rwmtx.Unlock()

	msg, exists := m.message[msgID]
	if !exists {
		return errors.New("message not found")
	}
	if msg.UserID != userID {
		return errors.New("you are not the author of this message")
	}

	delete(m.message, msgID)
	return nil
}

func (m *messageMemory) DeleteByRoom(ctx context.Context, roomID, userID int) error {
	m.rwmtx.Lock()
	defer m.rwmtx.Unlock()

	deleted := false
	for key, msg := range m.message {
		if msg.RoomID == roomID {
			delete(m.message, key)
			deleted = true
		}
	}

	if !deleted {
		return errors.New("no messages found in this room")
	}
	return nil
}
