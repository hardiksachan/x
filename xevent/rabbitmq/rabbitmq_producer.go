package rabbitmq

import (
	"context"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xevent"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitProducer is a wrapper around RabbitClient that exposes capabilities to send events
type RabbitProducer struct {
	client   *RabbitClient
	exchange xevent.Exchange
}

// NewRabbitProducer will create a new RabbitProducer
func NewRabbitProducer(client *RabbitClient, exchange xevent.Exchange) RabbitProducer {
	return RabbitProducer{
		client:   client,
		exchange: exchange,
	}
}

// Send will send a payload to the exchange
func (rp *RabbitProducer) Send(ctx context.Context, topic xevent.Topic, event *xevent.Event) error {
	op := xerrors.Op("xevent.RabbitProducer.Send")

	// TODO: make it traceable
	//nolint:exhaustruct
	publishing := amqp.Publishing{
		Type: event.Type,
		Body: event.Data,
	}

	err := rp.client.Send(ctx, string(rp.exchange), string(topic), publishing)
	if err != nil {
		return xerrors.E(op, err)
	}
	return nil
}
