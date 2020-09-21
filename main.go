package main

import (
	"fmt"
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
	mReceiver.Consume("notifications", helpers.Notification{}, func(msg interface{}) {
		fmt.Println("consuming notification")
		fmt.Println(msg)
		firebase, err := (&helpers.Firebase{}).InitApp()
		if err != nil {
			log.Fatalf("An error occurred while initializing firebase app : %s", err)
		}

		notification := msg.(helpers.Notification)
		for _, recipient := range notification.GetRecipients() {
			firebase.PushNotification(&notification, recipient)
		}
	})
}
