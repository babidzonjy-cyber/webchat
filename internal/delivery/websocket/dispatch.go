package websocket

import (
	"context"
	"fmt"
	"web-chat/internal/domain"
	"web-chat/internal/hub"
	"web-chat/internal/service"
)

func readWsResponse(ctx context.Context, msgSvc service.MessageService, userSvc service.UserService, msg *domain.Message, client *hub.Client) (*OutGoingMessage, error) {
	if err := msgSvc.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to save message %w", err)
	}

	username := "unknown"
	user, err := userSvc.GetByID(ctx, client.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for broadcast %w", err)
	}

	if user != nil {
		username = user.FullName
	}

	response := &OutGoingMessage{
		Type:      "message",
		ID:        msg.ID,
		UserID:    user.ID,
		Username:  username,
		Content:   msg.Text,
		CreatedAt: msg.CreatedAt,
	}

	return response, nil
}

func newWsMessage(client *hub.Client, incoming IncomingMessage) *domain.Message {
	return &domain.Message{
		Text:   incoming.Text,
		RoomID: client.RoomID,
		UserID: client.UserID,
	}
}
