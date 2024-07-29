package cache

import (
	"context"
	"time"
)

// Client defines the interface for connecting to the cache service
type Client interface {
	// Set Redis `SET key value [expiration]` command.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
	Get(ctx context.Context, key string) (string, error)

	// Get Redis `DELETE key` command. It returns redis.Nil error when key does not exist.
	Delete(ctx context.Context, key string) error

	// Check if key exist
	Exist(ctx context.Context, key string) (bool, error)

	// Close closes the client, releasing any open resources.
	Close(ctx context.Context, key string) error
}
