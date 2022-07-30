package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewProducer(queueName string) RabbitMQProducer {
	return &RabbitMQProducerImpl{
		rabbitMQClient: NewRabbitMQClient(),
		messageChannel: make(chan []byte),
		queueName:      queueName,
		workerCount:    1,
	}
}

type RabbitMQProducer interface {
	Send(body []byte)
	Start()
	Close()
}

type RabbitMQProducerImpl struct {
	rabbitMQClient RabbitMQClient
	messageChannel chan []byte
	queueName      string
	workerCount    int
}

func (producer *RabbitMQProducerImpl) Send(body []byte) {
	producer.messageChannel <- body
}

func (producer *RabbitMQProducerImpl) Start() {
	for i := 0; i < producer.workerCount; i++ {
		go producer.execute()
	}
}

func (producer *RabbitMQProducerImpl) execute() {
	channel, err := producer.rabbitMQClient.NewChannel()
	if err != nil {
		fmt.Sprintln("Error creating RabbitMQ channel for go routine")
		return
	}
	defer channel.Close()

	for {
		email := <-producer.messageChannel
		if channel.IsClosed() {
			return
		}
		if err := producer.publish(email, channel); err != nil {
			fmt.Sprintf("Error enque email: %s\n", email)
		}
	}
}

func (producer *RabbitMQProducerImpl) publish(body []byte, channel *amqp.Channel) error {
	q, err := producer.rabbitMQClient.DeclareQueue(channel, producer.queueName)
	if err != nil {
		return err
	}

	if err := producer.rabbitMQClient.Publish(channel, q, body); err != nil {
		return err
	}
	fmt.Printf(" [x] Sent %s\n", body)
	return nil
}

func (producer *RabbitMQProducerImpl) Close() {
	producer.rabbitMQClient.CloseConnections()
}
