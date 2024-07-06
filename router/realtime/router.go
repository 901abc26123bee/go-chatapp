package realtime_router

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	realtime_service "gsm/api/realtime"
	cors_middleware "gsm/middleware/cors"
	errors_middleware "gsm/middleware/errors"
	timeout_middleware "gsm/middleware/timeout"
	rediscache "gsm/pkg/cache/redis"
	gormpsql "gsm/pkg/orm/gorm"
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
	service realtime_service.RealtimeController
}

// NewRouter initialize routing information with controllers.
func NewRouter(ctx context.Context, config RouterConfig) (*RealtimeRouter, error) {
	// initialize orm with config.
	db, err := gormpsql.Initialize(config.SqlConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %v", err)
	}

	redisClient, err := rediscache.InitializeRedis(config.RedisConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %v", err)
	}

	realtimeController, err := realtime_service.NewRealtimeController(ctx, db, redisClient)
	if err != nil {
		return nil, fmt.Errorf("failed to new realtime controller: %v", err)
	}

	// Start listening for incoming chat messages
	go realtime_service.HandleMessages()

	r := gin.Default()

	// TODO: do not allow *
	corsHandler := cors_middleware.CorsHandler("*")
	errorHandler := errors_middleware.ErrorHandler()
	timeoutHandler := timeout_middleware.RequestTimeoutHandler()

	realtimeGroup := r.Group(path.Join("/api/realtime", realtimeVersion))
	realtimeGroup.Use(corsHandler, errorHandler, timeoutHandler)
	{
		realtimeGroup.GET("/healthz", getHealthz)
		{
			realtimeGroup.GET("/echo", realtimeController.Echo)
			realtimeGroup.GET("/ws", realtimeController.HandleWebsocketIO)
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
