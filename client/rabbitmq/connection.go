package rabbitmq

import (
	"fmt"
	"sync/atomic"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (client *RabbitMQClientImpl) GetConnection() (*amqp.Connection, error) {
	if err := client.addConnection(); err != nil {
		return nil, err
	}

	select {
	case conn := <-client.connections:
		if conn.IsClosed() {
			atomic.AddInt32(&client.initializedConnections, -1)
			return client.GetConnection()
		}
		return conn, nil
	}
}

func (client *RabbitMQClientImpl) ReturnConnection(connection *amqp.Connection) {
	client.connections <- connection
}

func (client *RabbitMQClientImpl) CloseConnections() {
	for {
		select {
		case connection := <-client.connections:
			if err := connection.Close(); err != nil {
				fmt.Sprintf("Error while cosing rabbitmq connection: %s\n", err.Error())
			}
		}
	}
}

func (client *RabbitMQClientImpl) addConnection() error {
	if client.initializedConnections < client.maxConnectionPoolSize {
		client.mutex.Lock()
		{
			if client.initializedConnections < client.maxConnectionPoolSize {
				conn, err := client.newConnection()
				if err != nil {
					return err
				}
				client.connections <- conn
				atomic.AddInt32(&client.initializedConnections, 1)
			}
		}
		client.mutex.Unlock()
	}
	return nil
}

func (client *RabbitMQClientImpl) newConnection() (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		client.username, client.password, client.host, client.port))
}
