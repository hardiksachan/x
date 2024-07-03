// Package outbox provides a transactional outbox
package outbox

import (
	"context"

	"github.com/hardiksachan/x/xerrors"
	"github.com/hardiksachan/x/xmessage"
	"github.com/hardiksachan/x/xretry"
)

const (
	failedMessagesChanSize = 10
)

// DataStore is the interface that wraps the Read method
type DataStore interface {
	GetUnsentPublishings(ctx context.Context) (<-chan *xmessage.Publishing, error)
	SetAsProcessed(ctx context.Context, id string) error
}

// PublishingStream is used to send publishing to queue
type PublishingStream interface {
	Send(context.Context, *xmessage.Publishing) error
}

// FailedPublishing is a publishing that failed to be dispatched
type FailedPublishing struct {
	Publishing *xmessage.Publishing
	Err        error
}

// Outbox is a transactional outbox
type Outbox struct {
	ds DataStore
	ps PublishingStream
	r  *xretry.Retrier
	fm chan *FailedPublishing
}

// New creates a new Outbox
func New(ds DataStore, p PublishingStream, r *xretry.Retrier) Outbox {
	return Outbox{
		ds: ds,
		ps: p,
		r:  r,
		fm: make(chan *FailedPublishing, failedMessagesChanSize),
	}
}

// Start will start the outbox
func (o *Outbox) Start(ctx context.Context) error {
	op := xerrors.Op("outbox.Outbox.Start")

	messages, err := o.ds.GetUnsentPublishings(ctx)
	if err != nil {
		return xerrors.E(op, err)
	}

	go o.startDispatcher(ctx, messages)

	return nil
}

func (o *Outbox) startDispatcher(ctx context.Context, publishings <-chan *xmessage.Publishing) {
	op := xerrors.Op("outbox.Outbox.startDispatcher")

	for {
		select {
		case p := <-publishings:
			err := o.r.Retry((func() error {
				return o.ps.Send(ctx, p)
			}))
			if err != nil {
				fm := &FailedPublishing{
					Publishing: p,
					Err:        xerrors.E(op, err),
				}

				o.fm <- fm
			}
			_ = o.ds.SetAsProcessed(ctx, p.Message.ID)
		case <-ctx.Done():
			return
		}
	}
}

// FailedPublishings returns a channel of failed publishings
func (o *Outbox) FailedPublishings() <-chan *FailedPublishing {
	return o.fm
}
