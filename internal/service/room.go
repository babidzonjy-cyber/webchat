package service

import (
	"context"
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
		return err
	}

	return nil
}

func (r *roomMemory) GetByID(ctx context.Context, id int) (*domain.Room, error) {
	room, err := r.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return room, err
}

func (r *roomMemory) GetAll(ctx context.Context) ([]*domain.Room, error) {
	rooms, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *roomMemory) Update(ctx context.Context, room *domain.Room) error {
	if err := r.repo.Update(ctx, room); err != nil {
		return err
	}

	return nil
}

func (r *roomMemory) Delete(ctx context.Context, id int) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
