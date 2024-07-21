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

	// // StreamAdd message to the end of stream(queue)
	// StreamAdd(ctx context.Context, args StreamAddArgs) error

	// Close closes the client, releasing any open resources.
	Close(ctx context.Context, key string) error
}

// type StreamAddArgs struct {
// 	Stream     string
// 	NoMkStream bool
// 	MaxLen     int64

// 	Limit  int64
// 	ID     string
// 	Values interface{}
// }

// type StreamReadArgs struct {
// 	Stream   []string
// 	Count   int64
// 	Block   time.Duration
// }

// type StreamReadMessages struct {
// 	Val []XStream
// }

// type XStream struct {
// 	Stream   string
// 	Messages []XMessage
// }

// type XMessage struct {
// 	ID     string
// 	Values map[string]interface{}
// }
