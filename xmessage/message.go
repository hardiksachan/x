// Package xmessage provides a asynchronous messaging system,
// With out-of-the-box support for transactional outbox
// and idempotent consumer.
package xmessage

// Topic is a topic of events.
type Topic string

// Message is a message that is published to a topic
type Message struct {
	// ID is the unique identifier of the message
	ID string
	// Type is the type of the message
	Type string
	// Payload is the message payload
	Payload []byte
}

// Publishing is a publishing of an event.
type Publishing struct {
	// Topic is the topic of the event.
	Topic Topic
	// Message is the message of the event.
	Message *Message
}

// Delivery is a delivery of an event.
type Delivery interface {
	Message() *Message
	Ack() error
	Nack() error
}
