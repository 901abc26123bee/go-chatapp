package realtime

import (
	"context"
	"fmt"

	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"gsm/pkg/stream"
)

// chatroomService defines the implementation of ChatRoomService interface
type chatroomService struct {
	redisClient  cache.Client
	streamClient stream.Client
}

// ChatRoomService defines the chat room service interface
type ChatRoomService interface {
	CreateChatRoom(ctx context.Context, req *CreateChatRoomRequest) error
	DeleteChatRoom(ctx context.Context) error
	QueryChatRoomHistory(ctx context.Context) error
}

// NewChatService init the chat service
func NewChatRoomService(redisClient cache.Client) ChatRoomService {
	return &chatroomService{
		redisClient: redisClient,
	}
}

type CreateChatRoomRequest struct {
	ChatRoomID string `json:"room_id"`
}

func (impl *chatroomService) CreateChatRoom(ctx context.Context, req *CreateChatRoomRequest) error {
	// TODO: create chatroom in db

	// create stream topic for chatroom when chat room is created for simplicity
	topicID := fmt.Sprintf("streamTopic:%s", req.ChatRoomID)
	_, err := impl.streamClient.CreateTopic(ctx, topicID)
	if err != nil {
		return errors.Errorf("failed to create stream topic: %v", err)
	}

	return nil
}

func (impl *chatroomService) DeleteChatRoom(ctx context.Context) error {
	return nil
}

func (impl *chatroomService) QueryChatRoomHistory(ctx context.Context) error {
	return nil
}
