package rediscache

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"

	"gsm/pkg/cache"
)

// InitializeRedis initial redis client
func InitializeRedis(redisConfigPath string) (cache.Client, error) {
	redisConfig, err := os.ReadFile(redisConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get db config: %v", err.Error())
	}
	opts, err := redis.ParseURL(string(redisConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to init redis: %v", err.Error())
	}

	client := redis.NewClient(opts)
	return Wrap(client), nil
}
