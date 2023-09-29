// Package outbox provides a transactional outbox
package outbox

import (
	"context"

	"github.com/Logistics-Coordinators/x/xerrors"
	"github.com/Logistics-Coordinators/x/xretry"
)

const (
	failedMessagesChanSize = 10
)

// DataStore is the interface that wraps the Read method
type DataStore interface {
	GetUnsentMessages(ctx context.Context) (<-chan *Message, error)
	SetAsProcessed(ctx context.Context, id string) error
}

// EventStream is used to send messages to the message broker
type EventStream interface {
	Send(*Message) error
}

// FailedMessage is a message that failed to be dispatched
type FailedMessage struct {
	Msg *Message
	Err error
}

// Outbox is a transactional outbox
type Outbox struct {
	ds DataStore
	es EventStream
	r  *xretry.Retrier
	fm chan *FailedMessage
}

// New creates a new Outbox
func New(ds DataStore, p EventStream, r *xretry.Retrier) Outbox {
	return Outbox{
		ds: ds,
		es: p,
		r:  r,
		fm: make(chan *FailedMessage, failedMessagesChanSize),
	}
}

// Start will start the outbox
func (o *Outbox) Start(ctx context.Context) error {
	op := xerrors.Op("outbox.Outbox.Start")

	messages, err := o.ds.GetUnsentMessages(ctx)
	if err != nil {
		return xerrors.E(op, err)
	}

	go o.StartDispatcher(ctx, messages)

	return nil
}

// StartDispatcher will start dispatching all messages in the xdispatcher.DataStore to xdispatcher.EventStream
func (o *Outbox) StartDispatcher(ctx context.Context, messages <-chan *Message) {
	op := xerrors.Op("outbox.Outbox.StartDispatcher")

	for {
		select {
		case m := <-messages:
			err := o.r.Retry((func() error {
				return o.es.Send(m)
			}))
			if err != nil {
				fm := &FailedMessage{
					Msg: m,
					Err: xerrors.E(op, err),
				}

				o.fm <- fm
			}
			_ = o.ds.SetAsProcessed(ctx, m.ID)
		case <-ctx.Done():
			return
		}
	}
}

// FailedMessages returns a channel of failed messages
func (o *Outbox) FailedMessages() <-chan *FailedMessage {
	return o.fm
}
