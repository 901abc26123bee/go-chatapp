package realtime

import (
	"github.com/gin-gonic/gin"

	"gsm/middleware/jwt"
	"gsm/pkg/cache"
	"gsm/pkg/errors"
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
	// get id from gin context parsed in auth middleware
	jwtClaimsID, ok := ctx.Get(jwt.JWTClaimID)
	if !ok {
		ctx.Error(errors.NewError(errors.InternalServerError, "failed to get id from access token"))
		return
	}
	userID, ok := jwtClaimsID.(string)
	if !ok {
		ctx.Error(errors.NewError(errors.InternalServerError, "failed to convert jwt claim id to string"))
		return
	}
	impl.connectService.HandleWebSocketStreamConnect(userID, ctx.Request, ctx.Writer)
}

func (impl *realtimeController) CreateChatRoom(ctx *gin.Context) {

}
