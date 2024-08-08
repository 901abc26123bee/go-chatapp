package realtime

import (
	"context"
	"fmt"

	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"gsm/pkg/stream"
)

// streamChatService defines the implementation of ChatService interface
type streamChatService struct {
	redisClient  cache.Client
	streamClient stream.Client
}

// StreamChatService defines the chat service interface
type StreamChatService interface {
	CreateChatRoom(ctx context.Context, req *CreateChatRoomRequest) error
	DeleteChatRoom(ctx context.Context) error
	QueryChatRoomHistory(ctx context.Context) error
}

// NewChatService init the chat service
func NewStreamChatService(redisClient cache.Client) StreamChatService {
	return &streamChatService{
		redisClient: redisClient,
	}
}

type CreateChatRoomRequest struct {
	ChatRoomID string `json:"room_id"`
}

func (impl *streamChatService) CreateChatRoom(ctx context.Context, req *CreateChatRoomRequest) error {
	// TODO: create chatroom in db

	// create stream topic for chatroom when chat room is created for simplicity
	topicID := fmt.Sprintf("streamTopic:%s", req.ChatRoomID)
	_, err := impl.streamClient.CreateTopic(ctx, topicID)
	if err != nil {
		return errors.Errorf("failed to create stream topic: %v", err)
	}

	return nil
}

func (impl *streamChatService) DeleteChatRoom(ctx context.Context) error {
	return nil
}

func (impl *streamChatService) QueryChatRoomHistory(ctx context.Context) error {
	return nil
}
