package streamredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisAcker struct {
	client     *redis.Client
	topicID    string
	xGroup     string
	messageIDs []string
}

func (acker *redisAcker) Ack(ctx context.Context) error {
	res := acker.client.XAck(ctx, acker.topicID, acker.xGroup, acker.messageIDs...)
	return res.Err()
}
