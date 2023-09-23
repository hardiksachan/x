package xevent

import "context"

type NoopBackbone struct{}

func (n *NoopBackbone) Send(context.Context, Topic, *Event) error {
	return nil
}

func (n *NoopBackbone) Listen() (<-chan *Delivery, error) {
	return nil, nil
}

func NewNoopBackbone() *NoopBackbone {
	return &NoopBackbone{}
}
