// Package queue provides the event system
package queue

import (
	"context"

	"github.com/Logistics-Coordinators/x/xmessage"
)

// Message is a notification of something that happened in the system.
type Message struct {
	ID   string
	Type string
	Data []byte
}

// Delivery is a delivery of an event.
type Delivery interface {
	Message() *Message
	Ack() error
	Nack() error
}

// Publisher is a producer of events.
type Publisher interface {
	// Send sends an event to the server.
	Send(ctx context.Context, topic xmessage.Topic, event *Message) error
}

// Consumer is a consumer of events.
type Consumer interface {
	// Listen listens for events from the server.
	Listen() (<-chan Delivery, error)
}
