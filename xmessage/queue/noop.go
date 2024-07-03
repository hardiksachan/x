package queue

import (
	"context"

	"github.com/hardiksachan/x/xmessage"
)

// NoopBackbone is a no-op implementation of Backbone.
type NoopBackbone struct{}

// Send implements queue.Backbone.
func (n *NoopBackbone) Send(context.Context, *xmessage.Publishing) error {
	return nil
}

// Listen implements queue.Backbone.
func (n *NoopBackbone) Listen() (<-chan xmessage.Delivery, error) {
	return nil, nil
}

// NewNoopBackbone returns a new NoopBackbone.
func NewNoopBackbone() *NoopBackbone {
	return &NoopBackbone{}
}
