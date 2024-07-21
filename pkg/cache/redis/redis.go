package rediscache

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"

	"gsm/pkg/cache"
)

// InitializeRedis initial redis client
func InitializeRedis(redisConfigPath string) (*redis.Client, error) {
	redisConfig, err := os.ReadFile(redisConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis config: %v", err.Error())
	}
	opts, err := redis.ParseURL(string(redisConfig))
	opts.PoolSize = 10

	if err != nil {
		return nil, fmt.Errorf("failed to init redis: %v", err.Error())
	}

	client := redis.NewClient(opts)
	return client, nil
}

// RedisWithCacheWrapper return redis client with cache interface
func RedisWithCacheWrapper(client *redis.Client) cache.Client {
	return Wrap(client)
}
