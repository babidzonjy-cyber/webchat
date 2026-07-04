package websocket

import (
	"context"
	"fmt"
	"web-chat/internal/domain"
	"web-chat/internal/hub"
	"web-chat/internal/service"
)

func handleIncomingMessage(ctx context.Context, msgSvc service.MessageService, userSvc service.UserService, msg *domain.Message, client *hub.Client) (*OutgoingMessage, error) {
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

	response := &OutgoingMessage{
		Type:      "message",
		ID:        msg.ID,
		UserID:    user.ID,
		Username:  username,
		Content:   msg.Text,
		CreatedAt: msg.CreatedAt,
	}

	return response, nil
}

func buildDomainMessage(client *hub.Client, incoming IncomingMessage) *domain.Message {
	return &domain.Message{
		Text:   incoming.Content,
		RoomID: client.RoomID,
		UserID: client.UserID,
	}
}
