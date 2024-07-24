package stream

import "context"

const DefaultTopicChannelCapacity int64 = 64

// Topic is a reference to a PubSub topic.
type Topic interface {
	// Exists reports whether the topic exists on the server.
	Exists(ctx context.Context) (bool, error)
	// Send publishes msg to the topic asynchronously.
	Send(ctx context.Context, msg *Message) PublishResult
}

// A PublishResult holds the result from a call to Publish.
type PublishResult interface {
	// Get returns the error result of a Publish call.
	Get(ctx context.Context) error
}
