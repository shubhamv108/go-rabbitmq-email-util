package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (client *RabbitMQClientImpl) NewChannel() (*amqp.Channel, error) {
	connection, err := client.GetConnection()
	if err != nil {
		return nil, err
	}
	defer client.ReturnConnection(connection)

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (client *RabbitMQClientImpl) Publish(channel *amqp.Channel, queue amqp.Queue, body []byte) error {
	return channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func (client *RabbitMQClientImpl) Consume(channel *amqp.Channel, queue amqp.Queue, autoAck bool) (<-chan amqp.Delivery, error) {
	return channel.Consume(
		queue.Name, // queue
		"",         // consumer
		autoAck,    // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
}
