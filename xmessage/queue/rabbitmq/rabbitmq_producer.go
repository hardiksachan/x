package rabbitmq

import (
	"context"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xmessage/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitProducer is a wrapper around RabbitClient that exposes capabilities to send events
type RabbitProducer struct {
	client   *RabbitClient
	exchange queue.Exchange
}

// NewRabbitProducer will create a new RabbitProducer
func NewRabbitProducer(client *RabbitClient, exchange queue.Exchange) RabbitProducer {
	return RabbitProducer{
		client:   client,
		exchange: exchange,
	}
}

// Send will send a payload to the exchange
func (rp *RabbitProducer) Send(ctx context.Context, topic queue.Topic, message *queue.Message) error {
	op := xerrors.Op("queue.RabbitProducer.Send")

	// TODO: make it traceable
	//nolint:exhaustruct
	publishing := amqp.Publishing{
		Type:      message.Type,
		Body:      message.Data,
		MessageId: message.ID,
	}

	err := rp.client.Send(ctx, string(rp.exchange), string(topic), publishing)
	if err != nil {
		return xerrors.E(op, err)
	}
	return nil
}
