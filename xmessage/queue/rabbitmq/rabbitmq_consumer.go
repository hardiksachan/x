package rabbitmq

import (
	"github.com/hardiksachan/x/xerrors"
	"github.com/hardiksachan/x/xmessage"
)

// RabbitConsumer is a wrapper around RabbitClient that exposes capabilities to listen for events
type RabbitConsumer struct {
	client *RabbitClient
	queue  string
}

// NewRabbitConsumer will create a new RabbitConsumer
func NewRabbitConsumer(client *RabbitClient, queue string) RabbitConsumer {
	return RabbitConsumer{
		client: client,
		queue:  queue,
	}
}

// Listen will return a channel that can be used to listen for events
func (rc RabbitConsumer) Listen() (<-chan xmessage.Delivery, error) {
	op := xerrors.Op("queue.RabbitConsumer.Listen")

	// Create a new Channel that can be used to listen for events
	// This channel will be used to listen for events
	rabbitDeliveries, err := rc.client.Consume(rc.queue, rc.queue, false)
	if err != nil {
		return nil, xerrors.E(op, err)
	}

	// Create a new channel that will be used to return deliveries
	deliveries := make(chan xmessage.Delivery)

	// Create a new go routine that will listen for events
	go func() {
		for d := range rabbitDeliveries {
			deliveries <- &RabbitDelivery{
				delivery: &d,
			}
		}
	}()

	return deliveries, nil
}
