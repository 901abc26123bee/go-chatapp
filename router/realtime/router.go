package realtime_router

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	realtime "gsm/api/realtime"
	cors "gsm/middleware/cors"
	errors "gsm/middleware/errors"
	"gsm/middleware/jwt"
	timeout "gsm/middleware/timeout"
	rediscache "gsm/pkg/cache/redis"
	"gsm/pkg/stream/streamredis"
	"gsm/pkg/util/sonyflake"

	// gormpsql "gsm/pkg/orm/gorm"
	realtimeutil "gsm/pkg/realtime"
)

// version of realtime server
const realtimeVersion = "v1"

// RouterConfig defines configs for dataset router
type RouterConfig struct {
	SqlConfigPath     string
	RedisConfigPath   string
	MongodbConfigPath string
	JwtSecret         string
}

// RealtimeRouter defines a gin engine for realtime router.
type RealtimeRouter struct {
	*gin.Engine
	service realtime.RealtimeController
}

// NewRouter initialize routing information with controllers.
func NewRouter(config RouterConfig) (*RealtimeRouter, error) {
	// initialize orm with config.
	// db, err := gormpsql.Initialize(config.SqlConfigPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize db: %v", err)
	// }

	// initialize cache with config.
	redisClient, err := rediscache.InitializeRedis(config.RedisConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %v", err)
	}
	cache := rediscache.RedisWithCacheWrapper(redisClient)
	stream := streamredis.NewRedisMessageStreamClient(redisClient)
	idGenerator, err := sonyflake.NewSonyFlake()
	if err != nil {
		return nil, fmt.Errorf("failed to  new SonyFlake: %v", err)
	}

	realtimeController, err := realtime.NewRealtimeController(cache, stream, idGenerator)
	if err != nil {
		return nil, fmt.Errorf("failed to new realtime controller: %v", err)
	}

	// Start listening for incoming chat messages(for testing)
	go realtimeutil.Broadcaster()

	r := gin.Default()

	// TODO: do not allow *
	corsHandler := cors.CorsHandler("*")
	errorHandler := errors.ErrorHandler()
	timeoutHandler := timeout.RequestTimeoutHandler()
	jwtHandler := jwt.HandleHeaderAuthorization(config.JwtSecret)

	r.Use(corsHandler, jwtHandler, errorHandler, timeoutHandler)

	realtimeGroup := r.Group(path.Join("/api/realtime", realtimeVersion))
	{
		realtimeGroup.GET("/healthz", getHealthz)
		{
			chatroomGroup := realtimeGroup.Group("/chatroom")
			{
				chatroomGroup.GET("/ws", realtimeController.HandleMemoryWebsocketIO)
				chatroomGroup.GET("/stream", realtimeController.HandleWebSocketStreamConnect)
			}
		}
	}

	return &RealtimeRouter{
		Engine:  r,
		service: realtimeController,
	}, nil
}

func getHealthz(ctx *gin.Context) {
	ctx.String(http.StatusOK, "alive!")
}
