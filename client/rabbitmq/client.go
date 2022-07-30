package rabbitmq

import (
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	newConnection() (*amqp.Connection, error)
	addConnection() error
	GetConnection() (*amqp.Connection, error)
	ReturnConnection(connection *amqp.Connection)
	CloseConnections()
	NewChannel() (*amqp.Channel, error)
	DeclareQueue(channel *amqp.Channel, queueName string) (amqp.Queue, error)
	Publish(channel *amqp.Channel, queue amqp.Queue, body []byte) error
	Consume(channel *amqp.Channel, queue amqp.Queue, autoAck bool) (<-chan amqp.Delivery, error)
}

type RabbitMQClientImpl struct {
	host                   string
	port                   string
	username               string
	password               string
	connections            chan *amqp.Connection
	initializedConnections int32
	maxConnectionPoolSize  int32
	mutex                  sync.RWMutex
}

func NewRabbitMQClient() RabbitMQClient {
	return &RabbitMQClientImpl{
		host:                  "localhost",
		port:                  "5672",
		username:              "guest",
		password:              "guest",
		maxConnectionPoolSize: 1,
		connections:           make(chan *amqp.Connection, 1+1),
	}
}
