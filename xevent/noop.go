package xevent

import "context"

// NoopBackbone is a no-op implementation of Backbone.
type NoopBackbone struct{}

// Send implements xevent.Backbone.
func (n *NoopBackbone) Send(context.Context, Topic, *Event) error {
	return nil
}

// Listen implements xevent.Backbone.
func (n *NoopBackbone) Listen() (<-chan *Delivery, error) {
	return nil, nil
}

// NewNoopBackbone returns a new NoopBackbone.
func NewNoopBackbone() *NoopBackbone {
	return &NoopBackbone{}
}
