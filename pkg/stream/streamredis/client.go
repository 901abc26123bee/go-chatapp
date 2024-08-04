package streamredis

import (
	"context"

	"gsm/pkg/stream"

	"github.com/redis/go-redis/v9"
)

const REDIS_ERROR_NO_SUCH_KEY = "ERR no such key"

type redisClient struct {
	client *redis.Client
}

// NewRedisMessageStreamClient creates a redis stream client.
func NewRedisMessageStreamClient(redis *redis.Client) stream.Client {
	return &redisClient{client: redis}
}

func (c *redisClient) Topic(id string) stream.Topic {
	return &redisTopic{
		id:     id,
		maxLen: stream.DefaultTopicChannelCapacity,
		client: c.client,
	}
}

func (c *redisClient) CreateTopic(ctx context.Context, topicID string) (stream.Topic, error) {
	return c.Topic(topicID), nil
}

func (c *redisClient) DeleteTopic(ctx context.Context, topicID string) error {
	return c.client.Del(ctx, topicID).Err()
}

func (c *redisClient) Subscription(xGroupID string, cfg *stream.SubscriptionConfig) stream.Subscription {
	return &redisSubscription{
		client:      c.client,
		xGroupID:    xGroupID,
		topicID:     cfg.TopicID,
		readStartID: cfg.ReadStartID,
	}
}

func (c *redisClient) CreateSubscription(ctx context.Context, xGroupID string, cfg *stream.SubscriptionConfig) (stream.Subscription, error) {
	topicExist, err := cfg.Topic.Exists(ctx)
	if err != nil {
		return nil, err
	}

	var res *redis.StatusCmd
	if topicExist {
		res = c.client.XGroupCreate(ctx, cfg.TopicID, xGroupID, cfg.ReadStartID)
	} else {
		// "$": start reading new messages added to the stream after the group is created,
		// ignoring any messages that were already present in the stream.
		res = c.client.XGroupCreateMkStream(ctx, cfg.TopicID, xGroupID, "$")
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	return c.Subscription(xGroupID, cfg), nil
}

func (c *redisClient) DeleteSubscription(ctx context.Context, topicID, xGroupID string) error {
	return c.client.XGroupDestroy(ctx, topicID, xGroupID).Err()
}
