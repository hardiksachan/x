// Package xevent provides the event system
package xevent

import "context"

// Topic is a topic of events.
type Topic string

// Event is a notification of something that happened in the system.
type Event struct {
	Type string
	Data []byte
}

// Delivery is a delivery of an event.
type Delivery interface {
	Event() *Event
	Ack() error
	Nack() error
}

// Publisher is a producer of events.
type Publisher interface {
	// Send sends an event to the server.
	Send(ctx context.Context, topic Topic, event *Event) error
}

// Consumer is a consumer of events.
type Consumer interface {
	// Listen listens for events from the server.
	Listen() (<-chan Delivery, error)
}
