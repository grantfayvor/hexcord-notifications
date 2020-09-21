package messaging

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

//Manager struct
type Manager struct {
	connection *amqp.Connection
}

//ManagerModel generic useful object
type ManagerModel struct {
	queue   amqp.Queue
	channel *amqp.Channel
}

//InitiateConnection init method for manager
func (m *Manager) InitiateConnection() {
	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_CONN"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ server: %s", err)
	}
	m.connection = conn
}

//CloseConnection this method should ideally be deferred
func (m *Manager) CloseConnection() error {
	return m.connection.Close()
}

func (m *Manager) createChannel() (*amqp.Channel, error) {
	return m.connection.Channel()
}

//DeclareQueue method on the manager object
func (m *Manager) DeclareQueue(name string) (*ManagerModel, error) {
	ch, err := m.createChannel()
	if err != nil {
		return nil, err
	}
	// defer ch.Close()

	q, err := ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &ManagerModel{channel: ch, queue: q}, nil
}
