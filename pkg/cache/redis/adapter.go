package rediscache

import (
	"context"
	"gsm/pkg/cache"
	"time"

	"github.com/go-redis/redis/v8"
)

// adapter defines the implementation for gorm to implement DB interface.
type adapter struct {
	client *redis.Client
}

// Wrap wraps a gorm db to orm DB.
func Wrap(cache *redis.Client) cache.Client {
	return &adapter{
		client: cache,
	}
}

func (a *adapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return a.client.Set(ctx, key, value, expiration).Err()
}

func (a *adapter) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key).Result()
}

func (a *adapter) Delete(ctx context.Context, key string) error {
	return a.client.Del(ctx, key).Err()
}
