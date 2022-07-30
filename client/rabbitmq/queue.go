package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (client *RabbitMQClientImpl) DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error) {
	return channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}
