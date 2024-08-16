package chat

import (
	"context"
	"gsm/pkg/cache"
	"gsm/pkg/stream"
)

// chatService defines the implementation of ChatService interface
type chatService struct {
	redisClient cache.Client
}

// ChatService defines the chat service interface
type ChatService interface {
	CreateChatRoom(ctx context.Context, req *CreateChatRoomRequest) (*CreateChatRoomResponse, error)
	DeleteChatRoom(ctx context.Context, req *CreateChatRoomRequest) (*CreateChatRoomResponse, error)
	JoinChatRoom(ctx context.Context, userID, chatRoomID, topicID, subID string) (stream.Subscription, error)
	LeaveChatRoom(ctx context.Context, topicID string, subID string) error
	PushMessage(ctx context.Context, userID string) error
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

func (impl *chatService) DeleteChatRoom(ctx context.Context, req *CreateChatRoomRequest) (*CreateChatRoomResponse, error) {
	return nil, nil
}

func (impl *chatService) JoinChatRoom(ctx context.Context, userID, chatRoomID, topicID, subID string) (stream.Subscription, error) {
	return nil, nil
}

func (impl *chatService) LeaveChatRoom(ctx context.Context, topicID string, subID string) error {
	return nil
}

func (impl *chatService) PushMessage(ctx context.Context, userID string) error {
	return nil
}
