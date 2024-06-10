package account_router

import (
	"context"
	"fmt"
	"net/http"
	"path"

	account_service "social-media-project/api/account"
	cors_middleware "social-media-project/middleware"

	"github.com/gin-gonic/gin"
)

// version of realtime server
const accountVersion = "v1"

// RouterConfig defines configs for account router
type RouterConfig struct {
	SqlConfigPath string
}

// AccountRouter defines a gin engine for account router.
type AccountRouter struct {
	*gin.Engine
	service account_service.AccountController
}

// NewRouter initialize routing information with controllers.
func NewRouter(ctx context.Context, config RouterConfig) (*AccountRouter, error) {
	accountController, err := account_service.NewAccountController(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to new account controller: %v", err)
	}

	r := gin.Default()

	// TODO: do not allow *
	corsHandler := cors_middleware.CorsHandler("*")

	realtimeGroup := r.Group(path.Join("/api/account", accountVersion))
	realtimeGroup.Use(corsHandler)
	{
		realtimeGroup.GET("/healthz", getHealthz)
		{
			realtimeGroup.GET("/user", accountController.GetUser)
			realtimeGroup.POST("/user", accountController.Register)
			realtimeGroup.POST("/login", accountController.Login)
		}
	}

	return &AccountRouter{
		Engine: r,
		service: accountController,
	}, nil
}

func getHealthz(ctx *gin.Context) {
	ctx.String(http.StatusOK, "alive!")
}
