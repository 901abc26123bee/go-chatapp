package streamredis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"gsm/pkg/messagebroker"
)

type redisTopic struct {
	id     string // topic channel id
	maxLen int64
	client *redis.Client
}

func (t *redisTopic) Exists(ctx context.Context) (bool, error) {
	_, err := t.client.XInfoStream(ctx, t.id).Result()
	if err == redis.Nil {
		// stream does not exist
		return false, nil
	} else if err != nil {
		// some other error occurred
		return false, err
	}

	return true, nil
}

func (t *redisTopic) Publish(ctx context.Context, msg *messagebroker.Message) messagebroker.PublishResult {
	rArg := &redis.XAddArgs{
		Stream: t.id,
		MaxLen: t.maxLen,
		ID:     msg.ID,
		Values: msg.Data,
	}
	res := t.client.XAdd(ctx, rArg)
	// TODO: check if exceed maxLen

	return &redisPublishResult{err: res.Err()}
}

type redisPublishResult struct {
	err error
}

func (pr *redisPublishResult) Get(ctx context.Context) error {
	return pr.err
}
