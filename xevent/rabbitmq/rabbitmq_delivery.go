package rabbitmq

import (
	"github.com/Logistics-Coordinators/x/xevent"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitDelivery is a wrapper around amqp.Delivery
type RabbitDelivery struct {
	delivery *amqp.Delivery
}

// Event will return the event
func (rd *RabbitDelivery) Event() *xevent.Event {
	return &xevent.Event{
		Type: rd.delivery.Type,
		Data: rd.delivery.Body,
	}
}

// Ack will acknowledge the delivery
func (rd *RabbitDelivery) Ack() error {
	return rd.delivery.Ack(false)
}

// Nack will negatively acknowledge the delivery
func (rd *RabbitDelivery) Nack() error {
	return rd.delivery.Nack(false, true)
}
