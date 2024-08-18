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

func (pr *redisPublishResult) Err(ctx context.Context) error {
	return pr.err
}

func (s *redisTopic) GetSubScriptions(ctx context.Context, prefix string) ([]string, error) {
	res := []string{}

	// scan for keys with the given prefix
	iter := s.client.Scan(ctx, 0, prefix, 0).Iterator()
	for iter.Next(ctx) {
		streamKey := iter.Val()

		// get groups for this stream
		groupListInfos, err := s.client.XInfoGroups(ctx, streamKey).Result()
		if err != nil {
			errMsg := fmt.Sprintf("%v", err)
			if errMsg == REDIS_ERROR_NO_SUCH_KEY {
				// skip if the stream does not exist
				continue
			}
			// some other error occurred
			return nil, err
		}

		for _, groupInfo := range groupListInfos {
			// Append each group name to the result list if the group belongs to the stream
			if groupInfo.Name != "" {
				res = append(res, groupInfo.Name)
			}
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
