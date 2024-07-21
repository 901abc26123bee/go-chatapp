package streamredis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"gsm/pkg/messagebroker"
)

type redisSubscription struct {
	client      *redis.Client
	xGroupID    string
	topicID     string
	consumerID  string
	readStartID string
}

func (s *redisSubscription) Exists(ctx context.Context) (bool, error) {
	_, err := s.client.XInfoGroups(ctx, s.xGroupID).Result()
	if err == redis.Nil {
		// stream does not exist
		return false, nil
	} else if err != nil {
		// some other error occurred
		return false, err
	}

	return true, nil
}

func (s *redisSubscription) Receive(ctx context.Context, f func(context.Context, *messagebroker.Message)) error {
	rArg := &redis.XReadGroupArgs{
		Group:    s.xGroupID,
		Consumer: s.consumerID,
		Streams:  []string{s.topicID},
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
		for _, stream := range res {
			for _, message := range stream.Messages {
				acker.messageIDs = append(acker.messageIDs, message.ID)
				msg := &messagebroker.Message{
					Acker:      acker,
					ID:         message.ID,
					Data:       []byte(fmt.Sprintf("%v", message.Values)),
					Attributes: message.Values,
				}
				// handle and ack the message.
				f(ctx, msg)
				msg.Acker.Ack(ctx)
			}
		}
	}
}
