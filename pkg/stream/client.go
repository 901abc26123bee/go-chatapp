package stream

import "context"

// Client is a messagebroker client
type Client interface {
	// Topic creates a reference to a topic in the client's project.
	Topic(topicID string) Topic
	// CreateTopic creates a new topic.
	CreateTopic(ctx context.Context, topicID string) (Topic, error)
	// DeleteTopic deletes a topic.
	DeleteTopic(ctx context.Context, topicID string) error

	// Subscription creates a reference to a subscription.
	Subscription(xGroupID string, cfg *SubscriptionConfig) Subscription
	// CreateSubscription creates a new subscription on a topic.
	CreateSubscription(ctx context.Context, xGroupID string, cfg *SubscriptionConfig) (Subscription, error)
	// DeleteSubscription deletes a subscription
	DeleteSubscription(ctx context.Context, xGroupID string) error
}

// SubscriptionConfig describes the configuration of a subscription.
type SubscriptionConfig struct {
	Topic       Topic
	TopicID     string
	ReadStartID string
}
