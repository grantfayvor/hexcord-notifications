package messaging

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

//Sender struct
type Sender struct {
	Manager
}

//Publish method on the Sender object
func (s *Sender) Publish(queueName string, message interface{}) error {
	manager, err := s.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return manager.channel.Publish("", manager.queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
