package rediscache

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"

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

	// test the connection to ensure it works
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return client, nil
}

// RedisWithCacheWrapper return redis client with cache interface
func RedisWithCacheWrapper(client *redis.Client) cache.Client {
	return Wrap(client)
}
