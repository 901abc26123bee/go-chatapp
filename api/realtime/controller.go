package realtime

import (
	"github.com/gin-gonic/gin"

	"gsm/pkg/cache"
	"gsm/pkg/realtime"
	"gsm/pkg/stream"
	"gsm/pkg/util/sonyflake"
)

// RealtimeController is the interface for realtime api
type RealtimeController interface {
	HandleMemoryWebsocketIO(*gin.Context)
	HandleWebSocketStreamConnect(*gin.Context)
	CreateChatRoom(*gin.Context)
}

// realtimeController defines the implementation of RealtimeController interface
type realtimeController struct {
	redisClient       cache.Client
	streamClient      stream.Client
	idGenerator       sonyflake.IDGenerator
	connectService    ConnectService
	streamChatService StreamChatService
}

// NewRealtimeController creates a new realtime controller
func NewRealtimeController(redisClient cache.Client, streamClient stream.Client, idGenerator sonyflake.IDGenerator) (RealtimeController, error) {
	return &realtimeController{
		redisClient:       redisClient,
		streamClient:      streamClient,
		idGenerator:       idGenerator,
		connectService:    NewConnectService(redisClient, streamClient, idGenerator),
		streamChatService: NewStreamChatService(redisClient),
	}, nil
}

// HandleMemoryWebsocketIO test websocket io with in-serve memory
func (impl *realtimeController) HandleMemoryWebsocketIO(ctx *gin.Context) {
	realtime.ServeWS(ctx.Writer, ctx.Request)
}

// HandleWebsocketIO handle socket io for client
func (impl *realtimeController) HandleWebSocketStreamConnect(ctx *gin.Context) {
	impl.connectService.HandleWebSocketStreamConnect(ctx.Writer, ctx.Request)
}

func (impl *realtimeController) CreateChatRoom(ctx *gin.Context) {

}
