package rabbitmq

import (
	"context"

	"github.com/hardiksachan/x/xerrors"
	"github.com/hardiksachan/x/xlog"
	"github.com/hardiksachan/x/xmessage"
	"github.com/hardiksachan/x/xmessage/queue"
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
func (rp *RabbitProducer) Send(ctx context.Context, publishing *xmessage.Publishing) error {
	op := xerrors.Op("queue.RabbitProducer.Send")

	xlog.Debugf("sending message to exchange %s with topic %s: %+v", rp.exchange, publishing.Topic, publishing.Message)

	// TODO: make it traceable
	//nolint:exhaustruct
	rbPublishing := amqp.Publishing{
		Type:      publishing.Message.Type,
		Body:      publishing.Message.Payload,
		MessageId: publishing.Message.ID,
	}

	err := rp.client.Send(ctx, string(rp.exchange), string(publishing.Topic), rbPublishing)
	if err != nil {
		return xerrors.E(op, err)
	}
	return nil
}
