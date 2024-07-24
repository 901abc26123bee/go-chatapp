package realtime

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"gsm/pkg/cache"
	"gsm/pkg/errors"
)

// chatService defines the implementation of ChatService interface
type chatService struct {
	redisClient cache.Client
}

// ChatService defines the chat service interface
type ChatService interface {
	CreateChatRoom(ctx context.Context) error
	QueryChatRoomHistory(ctx context.Context) error
}

// NewChatService init the chat service
func NewChatService(redisClient cache.Client) ChatService {
	return &chatService{
		redisClient: redisClient,
	}
}

func (impl *chatService) CreateChatRoom(ctx context.Context) error {
	chatroomID := ulid.Make().String()
	key := fmt.Sprintf("chatroom:%s", chatroomID)
	// set expiration to 0 indicates the key will not expire
	if err := impl.redisClient.Set(ctx, key, true, 0); err != nil {
		return errors.Errorf("failed to set user online in redis: %v\n", err)
	}
	return nil
}

func (impl *chatService) QueryChatRoomHistory(ctx context.Context) error {
	return nil
}
