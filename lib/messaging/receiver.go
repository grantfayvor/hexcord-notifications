package messaging

import (
	"encoding/json"
	"log"
)

//Receiver struct
type Receiver struct {
	Manager
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//Consume method on the Receiver object
func (r *Receiver) Consume(queueName string, incomingMsg interface{}, callback func(msg interface{})) {
	manager, err := r.DeclareQueue(queueName)
	if err != nil {
		failOnError(err, "Error creating the queue")
	}

	msgs, err := manager.channel.Consume(manager.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		failOnError(err, "Error occurred while trying to consume the queued message")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			json.Unmarshal(d.Body, &incomingMsg) //There might be a more reasonable thing to do here than ignore the error from reading this json
			callback(incomingMsg)
		}
	}()

	log.Printf("_____ [*] Waiting for messages. _____")
	<-forever
}
