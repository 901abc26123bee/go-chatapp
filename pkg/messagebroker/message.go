package messagebroker

import "context"

// Acker acks the message
type Acker interface {
	// Ack indicates successful processing of a Message passed to the Subscriber.Receive callback.
	Ack(ctx context.Context) error
}

// Message represents a Pub/Sub message.
type Message struct {
	Acker      Acker
	ID         string
	Data       []byte
	Attributes map[string]interface{}
}
