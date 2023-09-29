package queue

import (
	"context"

	"github.com/Logistics-Coordinators/x/xmessage"
)

// NoopBackbone is a no-op implementation of Backbone.
type NoopBackbone struct{}

// Send implements queue.Backbone.
func (n *NoopBackbone) Send(context.Context, xmessage.Topic, *Message) error {
	return nil
}

// Listen implements queue.Backbone.
func (n *NoopBackbone) Listen() (<-chan *Delivery, error) {
	return nil, nil
}

// NewNoopBackbone returns a new NoopBackbone.
func NewNoopBackbone() *NoopBackbone {
	return &NoopBackbone{}
}
