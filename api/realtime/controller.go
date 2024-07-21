package realtime

import (
	"github.com/gin-gonic/gin"

	"gsm/pkg/cache"
	"gsm/pkg/realtime"
)

// RealtimeController is the interface for realtime api
type RealtimeController interface {
	TestWebsocketIO(*gin.Context)
	HandleWebSocketConnect(*gin.Context)
	PushMessage(*gin.Context)
}

// realtimeController defines the implementation of RealtimeController interface
type realtimeController struct {
	redisClient    cache.Client
	connectService ConnectService
	chatService    ChatService
}

// NewRealtimeController creates a new realtime controller
func NewRealtimeController(redisClient cache.Client) (RealtimeController, error) {
	return &realtimeController{
		redisClient:    redisClient,
		connectService: NewConnectService(redisClient),
		chatService:    NewChatService(redisClient),
	}, nil
}

// HandleWebsocketIO handle socket io for client
func (impl *realtimeController) TestWebsocketIO(ctx *gin.Context) {
	realtime.ServeWS(ctx.Writer, ctx.Request)
}

// HandleWebsocketIO handle socket io for client
func (impl *realtimeController) HandleWebSocketConnect(ctx *gin.Context) {
	realtime.ServeWS(ctx.Writer, ctx.Request)
}

func (impl *realtimeController) PushMessage(ctx *gin.Context) {

}
