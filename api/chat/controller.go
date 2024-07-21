package chat

import (
	"context"
	"gsm/pkg/cache"

	"github.com/gin-gonic/gin"
)

// ChatController is the interface for chat api
type ChatController interface {
	CreateChatRoom(ctx *gin.Context)
	DeleteChatRoom(ctx *gin.Context)
	JoinChatRoom(ctx *gin.Context)
}

// chatController defines the implementation of ChatController interface
type chatController struct {
	chatService ChatService
}

// NewChatController creates a new chat controller
func NewChatController(ctx context.Context, redisClient cache.Client) (ChatController, error) {
	return &chatController{chatService: NewChatService(redisClient)}, nil
}

type CreateChatRoomRequest struct {
	Name string
}

type CreateChatRoomResponse struct {
	// empty
}

// CreateChatRoom create a chat room
func (impl *chatController) CreateChatRoom(ctx *gin.Context) {

}

// DeleteChatRoom delete a chat room
func (impl *chatController) DeleteChatRoom(ctx *gin.Context) {

}

// JoinChatRoom join a chat room
func (impl *chatController) JoinChatRoom(ctx *gin.Context) {

}
