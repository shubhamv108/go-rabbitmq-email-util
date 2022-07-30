package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConsumer(queueName string, messageProcessor func([]byte) error) RabbitMQConsumer {
	return &RabbitMQConsumerImpl{
		rabbitMQClient:   NewRabbitMQClient(),
		QueueName:        queueName,
		workerCount:      1,
		messageProcessor: messageProcessor,
	}
}

type RabbitMQConsumer interface {
	ProcessMessage(message []byte) error
	Start()
	Close()
}

type RabbitMQConsumerImpl struct {
	rabbitMQClient   RabbitMQClient
	QueueName        string
	workerCount      int
	messageProcessor func([]byte) error
}

func (consumer *RabbitMQConsumerImpl) Start() {
	for i := 0; i < consumer.workerCount; i++ {
		/*go*/ consumer.execute()
	}
}

func (consumer *RabbitMQConsumerImpl) execute() {
	channel, err := consumer.rabbitMQClient.NewChannel()
	if err != nil {
		fmt.Sprintln("Error creating RabbitMQ channel for go routine")
	}
	defer channel.Close()

	for {
		if channel.IsClosed() {
			return
		}
		if err := consumer.poll(channel); err != nil {
			fmt.Sprintf("Error polling message from rabbitmq ( Queue name: %s, Error: %s )\n", consumer.QueueName, err.Error())
		}
	}
}

func (consumer *RabbitMQConsumerImpl) poll(channel *amqp.Channel) error {
	q, err := consumer.rabbitMQClient.DeclareQueue(channel, consumer.QueueName)
	if err != nil {
		return err
	}

	msgs, err := consumer.rabbitMQClient.Consume(channel, q, false)
	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			msg := d.Body
			fmt.Printf("Received a message: %s\n", d.Body)
			if err := consumer.ProcessMessage(msg); err != nil {
				fmt.Sprintf("Error processing polled message: %s\n", err.Error())
			} else {
				channel.Ack(d.DeliveryTag, false)
			}
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever

	return nil
}

func (consumer *RabbitMQConsumerImpl) ProcessMessage(message []byte) error {
	return consumer.messageProcessor(message)
}

func (consumer *RabbitMQConsumerImpl) Close() {
	consumer.rabbitMQClient.CloseConnections()
}
