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
	// Topic is the name of the topic to which the message is published
	Topic string
	// Type is the type of the message
	Type string
	// Payload is the message payload
	Payload []byte
}
