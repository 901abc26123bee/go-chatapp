package streamredis

import (
	"context"

	"gsm/pkg/messagebroker"

	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	client *redis.Client
}

// NewRedisMessageBrokerClient creates a redis broker client.
func NewRedisMessageBrokerClient(redis *redis.Client) messagebroker.Client {
	return &redisClient{client: redis}
}

func (c *redisClient) Topic(id string) messagebroker.Topic {
	return &redisTopic{
		id:     id,
		maxLen: messagebroker.DefaultTopicChannelCapacity,
		client: c.client,
	}
}

func (c *redisClient) CreateTopic(ctx context.Context, topicID string) (messagebroker.Topic, error) {
	return c.Topic(topicID), nil
}

func (c *redisClient) DeleteTopic(ctx context.Context, topicID string) error {
	return c.client.Del(ctx, topicID).Err()
}

func (c *redisClient) Subscription(xGroupID string, cfg *messagebroker.SubscriptionConfig) messagebroker.Subscription {
	return &redisSubscription{
		client:      c.client,
		xGroupID:    xGroupID,
		topicID:     cfg.TopicID,
		readStartID: cfg.ReadStartID,
	}
}

func (c *redisClient) CreateSubscription(ctx context.Context, xGroupID string, cfg *messagebroker.SubscriptionConfig) (messagebroker.Subscription, error) {
	res := c.client.XGroupCreate(ctx, xGroupID, cfg.TopicID, cfg.ReadStartID)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return c.Subscription(xGroupID, cfg), nil
}

func (c *redisClient) DeleteSubscription(ctx context.Context, xGroupID string) error {
	return nil
}
