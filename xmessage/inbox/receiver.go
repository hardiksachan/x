// Package inbox provides a idempotent inbox
package inbox

import (
	"context"

	"github.com/hardiksachan/x/xerrors"
	"github.com/hardiksachan/x/xmessage"
)

// Receiver is a idempotent inbox
type Receiver struct {
	r Repository
}

// NewReceiver creates a new inbox.Receiver
func NewReceiver(r Repository) *Receiver {
	return &Receiver{
		r: r,
	}
}

// Receive will receive a message
func (i *Receiver) Receive(ctx context.Context, deliveries <-chan xmessage.Delivery) error {
	go func() {
		for delivery := range deliveries {
			message := delivery.Message()
			err := i.r.SaveMessage(ctx, message)
			if err != nil {
				if xerrors.ErrorCode(err) != xerrors.Exists {
					_ = delivery.Nack()
					continue
				}
			}
			_ = delivery.Ack()
		}
	}()

	return nil
}
