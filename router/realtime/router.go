package realtime_router

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	realtime "gsm/api/realtime"
	cors "gsm/middleware/cors"
	errors "gsm/middleware/errors"
	timeout "gsm/middleware/timeout"
	rediscache "gsm/pkg/cache/redis"

	// gormpsql "gsm/pkg/orm/gorm"
	realtimeutil "gsm/pkg/realtime"
)

// version of realtime server
const realtimeVersion = "v1"

// RouterConfig defines configs for dataset router
type RouterConfig struct {
	SqlConfigPath   string
	RedisConfigPath string
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

	realtimeController, err := realtime.NewRealtimeController(cache)
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

	realtimeGroup := r.Group(path.Join("/api/realtime", realtimeVersion))
	realtimeGroup.Use(corsHandler, errorHandler, timeoutHandler)
	{
		realtimeGroup.GET("/healthz", getHealthz)
		{
			realtimeGroup.GET("/ws", realtimeController.TestWebsocketIO)
			realtimeGroup.GET("/connect", realtimeController.HandleWebSocketConnect)
			pushGroup := realtimeGroup.Group("/push")
			{
				pushGroup.GET("/message", realtimeController.HandleWebSocketConnect)
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
