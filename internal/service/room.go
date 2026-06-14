package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"web-chat/internal/domain"
)

type RoomService interface {
	Create(ctx context.Context, room *domain.Room) error
	GetByID(ctx context.Context, id int) (*domain.Room, error)
	GetAll(ctx context.Context) ([]*domain.Room, error)
	Update(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, id int) error
}

type roomMemory struct {
	rooms  map[int]*domain.Room
	nextID int
	rwmtx  sync.RWMutex
}

func NewRoomMemory() *roomMemory {
	return &roomMemory{
		rooms:  make(map[int]*domain.Room),
		nextID: 1,
	}
}

func (r *roomMemory) Create(ctx context.Context, room *domain.Room) error {
	r.rwmtx.Lock()
	defer r.rwmtx.Unlock()

	room.ID = r.nextID
	r.nextID++
	room.CreatedAt = time.Now()
	r.rooms[room.ID] = room

	return nil
}

func (r *roomMemory) GetByID(ctx context.Context, id int) (*domain.Room, error) {
	r.rwmtx.RLock()
	defer r.rwmtx.RUnlock()

	if val, exists := r.rooms[id]; exists {
		return val, nil
	}

	return nil, errors.New("there is no room with that id")
}

func (r *roomMemory) GetAll(ctx context.Context) ([]*domain.Room, error) {
	r.rwmtx.RLock()
	defer r.rwmtx.RUnlock()

	roomsSlice := make([]*domain.Room, 0, len(r.rooms))

	if len(r.rooms) == 0 {
		return roomsSlice, nil
	}

	for _, v := range r.rooms {
		roomsSlice = append(roomsSlice, v)
	}

	return roomsSlice, nil
}

func (r *roomMemory) Update(ctx context.Context, room *domain.Room) error {
	r.rwmtx.Lock()
	defer r.rwmtx.Unlock()

	if _, exists := r.rooms[room.ID]; !exists {
		return fmt.Errorf("room %d not found", room.ID)
	}

	r.rooms[room.ID] = room
	return nil
}

func (r *roomMemory) Delete(ctx context.Context, id int) error {
	r.rwmtx.Lock()
	defer r.rwmtx.Unlock()

	if _, exists := r.rooms[id]; exists {
		delete(r.rooms, id)
		return nil
	}

	return errors.New("there is no room with that id")
}
