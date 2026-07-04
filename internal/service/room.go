package service

import (
	"context"
	"fmt"
	"web-chat/internal/domain"
	"web-chat/internal/repository"
)

type RoomService interface {
	Create(ctx context.Context, room *domain.Room) error
	GetByID(ctx context.Context, id int) (*domain.Room, error)
	GetAll(ctx context.Context) ([]*domain.Room, error)
	Update(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, id int) error
}

type roomMemory struct {
	repo repository.RoomRepository
}

func NewRoomMemory(repo repository.RoomRepository) *roomMemory {
	return &roomMemory{
		repo: repo,
	}
}

func (r *roomMemory) Create(ctx context.Context, room *domain.Room) error {
	if err := r.repo.Create(ctx, room); err != nil {
		return fmt.Errorf("service.Create room: %d, error %w", room.ID, err)
	}

	return nil
}

func (r *roomMemory) GetByID(ctx context.Context, id int) (*domain.Room, error) {
	room, err := r.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("service.GetByID room %d: %w", id, err)
	}

	return room, nil
}

func (r *roomMemory) GetAll(ctx context.Context) ([]*domain.Room, error) {
	rooms, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetAll error: %w", err)
	}

	return rooms, nil
}

func (r *roomMemory) Update(ctx context.Context, room *domain.Room) error {
	if err := r.repo.Update(ctx, room); err != nil {
		return fmt.Errorf("service.Update room: %d, error: %w", room.ID, err)
	}

	return nil
}

func (r *roomMemory) Delete(ctx context.Context, id int) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service.Delete room %d: %w", id, err)
	}

	return nil
}
