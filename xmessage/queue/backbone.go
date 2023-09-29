// Package queue provides the event system
package queue

import (
	"context"

	"github.com/Logistics-Coordinators/x/xmessage"
)

// Delivery is a delivery of an event.
type Delivery interface {
	Message() *xmessage.Message
	Ack() error
	Nack() error
}

// Publisher is a producer of events.
type Publisher interface {
	// Send sends an event to the server.
	Send(ctx context.Context, publishing *xmessage.Publishing) error
}

// Consumer is a consumer of events.
type Consumer interface {
	// Listen listens for events from the server.
	Listen() (<-chan Delivery, error)
}
