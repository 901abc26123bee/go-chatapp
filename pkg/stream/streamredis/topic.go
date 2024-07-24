package streamredis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"gsm/pkg/stream"
)

type redisTopic struct {
	id     string // topic channel id
	maxLen int64
	client *redis.Client
}

func (t *redisTopic) Exists(ctx context.Context) (bool, error) {
	// TODO: find a better way
	res := t.client.XInfoStream(ctx, t.id)
	err := res.Err()
	errMsg := fmt.Sprintf("%v", err)
	if errMsg == REDIS_ERROR_NO_SUCH_KEY {
		return false, nil
	} else if err != nil {
		// some other error occurred
		return false, err
	}

	return true, nil
}

func (t *redisTopic) Send(ctx context.Context, msg *stream.Message) stream.PublishResult {
	rArg := &redis.XAddArgs{
		Stream: t.id,
		MaxLen: t.maxLen,
		ID:     msg.ID,
		Values: msg.Attributes,
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
