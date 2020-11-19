package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grantfayvor/hexcord-notifications/helpers"
	"github.com/grantfayvor/hexcord-notifications/lib"
	"github.com/grantfayvor/hexcord-notifications/lib/messaging"
	"github.com/grantfayvor/hexcord-notifications/lib/notification"
	notifier "github.com/grantfayvor/hexcord-notifications/lib/notification"

	"github.com/Kamva/mgm"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initDB()

	go func() {
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

			notification := &notification.Notification{}
			err = json.Unmarshal(received, notification)
			if err != nil {
				log.Printf("An error occurred while parsing the json to notification object : %s", err)
				return
			}

			notifier.SaveNotification(notification)

			for _, recipient := range notification.GetRecipients() {
				firebase.PushNotification(notification, recipient)
			}
		})
	}()

	lib.InitializeRoutes()

	fmt.Printf("Starting server at port %s\n", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func initDB() error {
	return mgm.SetDefaultConfig(nil, "screen_recorder",
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
}
