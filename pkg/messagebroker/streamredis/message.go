package streamredis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type redisAcker struct {
	client     *redis.Client
	topicID    string
	xGroup     string
	messageIDs []string
}

func (acker *redisAcker) Ack(ctx context.Context) error {
	res := acker.client.XAck(ctx, acker.topicID, acker.topicID, acker.messageIDs...)
	return res.Err()
}
