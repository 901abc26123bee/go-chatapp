package realtime

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"gsm/pkg/cache"
	"gsm/pkg/errors"
)

// streamChatService defines the implementation of ChatService interface
type streamChatService struct {
	redisClient cache.Client
}

// StreamChatService defines the chat service interface
type StreamChatService interface {
	CreateChatRoom(ctx context.Context) error
	QueryChatRoomHistory(ctx context.Context) error
}

// NewChatService init the chat service
func NewStreamChatService(redisClient cache.Client) StreamChatService {
	return &streamChatService{
		redisClient: redisClient,
	}
}

func (impl *streamChatService) CreateChatRoom(ctx context.Context) error {
	chatroomID := ulid.Make().String()
	key := fmt.Sprintf("chatroom:%s", chatroomID)
	// set expiration to 0 indicates the key will not expire
	if err := impl.redisClient.Set(ctx, key, true, 0); err != nil {
		return errors.Errorf("failed to set user online in redis: %v\n", err)
	}
	return nil
}

func (impl *streamChatService) QueryChatRoomHistory(ctx context.Context) error {
	return nil
}
