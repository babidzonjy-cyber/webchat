package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type OnlineRepository interface {
	AddOnline(roomID, userID int) error
	RemoveOnline(roomID, userID int) error
	GetOnlineCount(roomID int) (int, error)
	GetOnlineUsers(roomID int) ([]int, error)
	IsOnline(roomID, userID int) (bool, error)
}

type RedisOnline struct {
	client *redis.Client
}

func NewRedisOnline(addr string) (*RedisOnline, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connect: %w", err)
	}
	return &RedisOnline{client: client}, nil
}

func (ro *RedisOnline) key(roomID int) string {
	return fmt.Sprintf("online:room:%d", roomID)
}

func (ro *RedisOnline) AddOnline(roomID, userID int) error {
	return ro.client.SAdd(context.Background(), ro.key(roomID), userID).Err()
}

func (ro *RedisOnline) RemoveOnline(roomID, userID int) error {
	return ro.client.SRem(context.Background(), ro.key(roomID), userID).Err()
}

func (ro *RedisOnline) GetOnlineCount(roomID int) (int, error) {
	count, err := ro.client.SCard(context.Background(), ro.key(roomID)).Result()

	return int(count), err
}

func (ro *RedisOnline) GetOnlineUsers(roomID int) ([]int, error) {
	members, err := ro.client.SMembers(context.Background(), ro.key(roomID)).Result()

	if err != nil {
		return nil, err
	}

	users := make([]int, 0, len(members))
	for _, m := range members {
		id, err := strconv.Atoi(m)
		if err != nil {
			continue
		}
		users = append(users, id)
	}

	return users, nil
}

func (ro *RedisOnline) IsOnline(roomID, userID int) (bool, error) {
	return ro.client.SIsMember(context.Background(), ro.key(roomID), userID).Result()
}

func (ro *RedisOnline) Close() error {
	return ro.client.Close()
}
