package realtime

import (
	"context"
	"gsm/pkg/cache"
)

// chatService defines the implementation of ChatService interface
type chatService struct {
	redisClient cache.Client
}

// ChatService defines the chat service interface
type ChatService interface {
	PushMessage(ctx context.Context) error
}

// NewChatService init the chat service
func NewChatService(redisClient cache.Client) ChatService {
	return &chatService{
		redisClient: redisClient,
	}
}

func (impl *chatService) PushMessage(ctx context.Context) error {
	return nil
}
