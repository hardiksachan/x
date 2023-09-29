package xqueue

import "context"

// NoopBackbone is a no-op implementation of Backbone.
type NoopBackbone struct{}

// Send implements xqueue.Backbone.
func (n *NoopBackbone) Send(context.Context, Topic, *Message) error {
	return nil
}

// Listen implements xqueue.Backbone.
func (n *NoopBackbone) Listen() (<-chan *Delivery, error) {
	return nil, nil
}

// NewNoopBackbone returns a new NoopBackbone.
func NewNoopBackbone() *NoopBackbone {
	return &NoopBackbone{}
}
