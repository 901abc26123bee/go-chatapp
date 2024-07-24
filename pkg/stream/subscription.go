package stream

import "context"

// Subscription is the interface for the messagebroker subscription service.
type Subscription interface {
	// Exists reports whether the subscription exists on the server.
	Exists(ctx context.Context) (bool, error)
	// Receive calls fto handle messages receive from the subscription.
	Receive(ctx context.Context, f func(context.Context, *Message)) error
}
