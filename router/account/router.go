package account_router

import (
	"context"
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
	"gsm/pkg/mdb"
)

// version of realtime server
const accountVersion = "v1"

// RouterConfig defines configs for account router
type RouterConfig struct {
	SqlConfigPath     string
	DBKey             string
	RedisConfigPath   string
	MongodbConfigPath string
}

// AccountRouter defines a gin engine for account router.
type AccountRouter struct {
	*gin.Engine
	service account.AccountController
}

// NewRouter initialize routing information with controllers.
func NewRouter(config RouterConfig) (*AccountRouter, error) {
	ctx := context.Background()
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

	mongodbClient, err := mdb.InitializeWithEncryptedKey(ctx, config.MongodbConfigPath, config.DBKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mongodb: %v", err)
	}

	accountController, err := account.NewAccountController(cache, mongodbClient, config.DBKey)
	if err != nil {
		return nil, fmt.Errorf("failed to new account controller: %v", err)
	}

	r := gin.Default()
	// TODO: do not allow *
	corsHandler := cors.CorsHandler("*")
	errorHandler := errors.ErrorHandler()
	timeoutHandler := timeout.RequestTimeoutHandler()
	r.Use(corsHandler, errorHandler, timeoutHandler)

	accountGroup := r.Group(path.Join("/api/account", accountVersion))
	{
		accountGroup.GET("/healthz", getHealthz)
		{
			accountGroup.GET("/user", accountController.GetUser)
			accountGroup.POST("/user", accountController.CreateUser)
			accountGroup.POST("/login", accountController.Login)
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
