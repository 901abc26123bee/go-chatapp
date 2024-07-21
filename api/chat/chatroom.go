package chat

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
	CreateChatRoom(ctx context.Context, req *CreateChatRoomRequest) (*CreateChatRoomResponse, error)
}

// NewChatService init the chat service
func NewChatService(redisClient cache.Client) ChatService {
	return &chatService{
		redisClient: redisClient,
	}
}

func (impl *chatService) CreateChatRoom(ctx context.Context, req *CreateChatRoomRequest) (*CreateChatRoomResponse, error) {
	// create chat room in db

	// create chat room stream queue

	// involve user to chatroom

	return nil, nil
}
