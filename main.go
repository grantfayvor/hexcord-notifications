package main

import (
	"encoding/json"
	"log"

	"github.com/grantfayvor/hexcord-notifications/helpers"
	"github.com/grantfayvor/hexcord-notifications/lib/messaging"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mReceiver := &messaging.Receiver{}
	mReceiver.InitiateConnection()
	mReceiver.Consume("notifications", func(msg map[string]interface{}) {
		firebase, err := (&helpers.Firebase{}).InitApp()
		if err != nil {
			log.Fatalf("An error occurred while initializing firebase app : %s", err)
		}

		received, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("An error occurred while marshalling the message : %s", err)
		}

		notification := &helpers.Notification{}
		err = json.Unmarshal(received, notification)
		if err != nil {
			log.Fatalf("An error occurred while parsing the json to notification object : %s", err)
		}

		for _, recipient := range notification.GetRecipients() {
			firebase.PushNotification(notification, recipient)
		}
	})
}
