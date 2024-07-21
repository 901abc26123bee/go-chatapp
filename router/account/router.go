package account_router

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	account "gsm/api/account"
	cors "gsm/middleware/cors"
	errors "gsm/middleware/errors"
	timeout "gsm/middleware/timeout"

	// gormpsql "gsm/pkg/orm/gorm"
	rediscache "gsm/pkg/cache/redis"
)

// version of realtime server
const accountVersion = "v1"

// RouterConfig defines configs for account router
type RouterConfig struct {
	SqlConfigPath   string
	DBKey           string
	RedisConfigPath string
}

// AccountRouter defines a gin engine for account router.
type AccountRouter struct {
	*gin.Engine
	service account.AccountController
}

// NewRouter initialize routing information with controllers.
func NewRouter(config RouterConfig) (*AccountRouter, error) {
	// initialize orm with config.
	// db, err := gormpsql.InitializeWithEncryptedKey(config.SqlConfigPath, config.DBKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize db: %v", err)
	// }

	// initialize cache with config.
	redisClient, err := rediscache.InitializeRedis(config.RedisConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %v", err)
	}
	cache := rediscache.RedisWithCacheWrapper(redisClient)

	accountController, err := account.NewAccountController(cache)
	if err != nil {
		return nil, fmt.Errorf("failed to new account controller: %v", err)
	}

	r := gin.Default()

	// TODO: do not allow *
	corsHandler := cors.CorsHandler("*")
	errorHandler := errors.ErrorHandler()
	timeoutHandler := timeout.RequestTimeoutHandler()

	realtimeGroup := r.Group(path.Join("/api/account", accountVersion))
	realtimeGroup.Use(corsHandler, errorHandler, timeoutHandler)
	{
		realtimeGroup.GET("/healthz", getHealthz)
		{
			realtimeGroup.GET("/user", accountController.GetUser)
			realtimeGroup.POST("/user", accountController.CreateUser)
			realtimeGroup.POST("/login", accountController.Login)
		}
	}

	return &AccountRouter{
		Engine:  r,
		service: accountController,
	}, nil
}

func getHealthz(ctx *gin.Context) {
	ctx.String(http.StatusOK, "alive!")
}
