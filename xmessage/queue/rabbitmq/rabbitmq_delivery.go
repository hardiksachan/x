package rabbitmq

import (
	"github.com/Logistics-Coordinators/x/xmessage/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitDelivery is a wrapper around amqp.Delivery
type RabbitDelivery struct {
	delivery *amqp.Delivery
}

// Message will return the event
func (rd *RabbitDelivery) Message() *queue.Message {
	return &queue.Message{
		ID:   rd.delivery.MessageId,
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
