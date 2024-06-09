package realtime_router

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	realtime_service "social-media-project/api/realtime"
	cors_middleware "social-media-project/middleware"
)

// version of realtime server
const realtimeVersion = "v1"

// RouterConfig defines configs for dataset router
type RouterConfig struct {
	SqlConfigPath string
}

// RealtimeRouter defines a gin engine for realtime router.
type RealtimeRouter struct {
	*gin.Engine
	service realtime_service.RealtimeController
}

// NewRouter initialize routing information with controllers.
func NewRouter(ctx context.Context, config RouterConfig) (*RealtimeRouter, error) {
	realtimeController, err := realtime_service.NewRealtimeController(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to new realtime controller: %v", err)
	}

	r := gin.Default()

	// TODO: do not allow *
	corsHandler := cors_middleware.CorsHandler("*")

	realtimeGroup := r.Group(path.Join("/api/realtime", realtimeVersion))
	realtimeGroup.Use(corsHandler)
	{
		realtimeGroup.GET("/healthz", getHealthz)
		{
			realtimeGroup.GET("/echo", realtimeController.Echo)
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
