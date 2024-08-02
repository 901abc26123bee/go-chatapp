package streamredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"gsm/pkg/stream"
)

type redisSubscription struct {
	client      *redis.Client
	xGroupID    string
	topicID     string
	readStartID string
}

func (s *redisSubscription) Exists(ctx context.Context) (bool, error) {
	// TODO: find a better way
	_, err := s.client.XInfoGroups(ctx, s.xGroupID).Result()
	errMsg := fmt.Sprintf("%v", err)
	if errMsg == REDIS_ERROR_NO_SUCH_KEY {
		return false, nil
	} else if err != nil {
		// some other error occurred
		return false, err
	}

	return true, nil
}

func (s *redisSubscription) Receive(ctx context.Context, f func(context.Context, *stream.Message)) error {
	rArg := &redis.XReadGroupArgs{
		Group:    s.xGroupID,
		Consumer: s.xGroupID,               // set consumer to group since each group will only consume by distinct member
		Streams:  []string{s.topicID, ">"}, // only get messages that were added after the last acknowledgment or after the consumer group was created.
		Count:    16,
		Block:    300 * time.Second, // block for 5 min
		NoAck:    false,
	}

	for {
		res, err := s.client.XReadGroup(ctx, rArg).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		acker := &redisAcker{
			client:     s.client,
			topicID:    s.topicID,
			xGroup:     s.xGroupID,
			messageIDs: []string{},
		}
		for _, xStream := range res {
			for _, message := range xStream.Messages {
				acker.messageIDs = append(acker.messageIDs, message.ID)
				msg := &stream.Message{
					Acker:      acker,
					ID:         message.ID,
					Attributes: message.Values,
				}
				// handle and ack the message.
				f(ctx, msg)
				msg.Acker.Ack(ctx)
			}
		}
	}
}
