package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/getsentry/sentry-go"
	"github.com/grantfayvor/hexcord-notifications/helpers"
	"github.com/grantfayvor/hexcord-notifications/lib"
	"github.com/grantfayvor/hexcord-notifications/lib/mailing"
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

	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_CLIENT"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	if err := initDB(); err != nil {
		log.Fatalf("An error occurred while trying to initialize the DB : %s", err)
	}

	go func() {
		mReceiver := &messaging.Receiver{}
		mReceiver.InitiateConnection()
		mReceiver.Consume(os.Getenv("RABBIT_MQ_CONN_NOTIFICATION_QUEUE"), func(msg map[string]interface{}) {
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
				notificationMessage := notification.GetMessage()
				mailTemplate := strings.ReplaceAll(
					strings.Replace(
						strings.Replace(
							strings.Replace(
								strings.Replace(mailing.RecordingNotificationMailTemplate, "{{recipientName}}", notificationMessage["recipientName"], 1),
								"{{thumbNail}}",
								notificationMessage["recordThumbnail"],
								1,
							),
							"{{senderName}}",
							notificationMessage["senderName"],
							1,
						),
						"{{notificationMessage}}",
						notificationMessage["title"],
						1,
					),
					"{{recordingLink}}",
					notificationMessage["recordingLink"],
				)

				err := mailing.NewMailer().
					InitMessage(notificationMessage["recipientEmail"], notificationMessage["title"]).
					SendMail(mailTemplate)
				fmt.Println("==================")
				fmt.Println(err)
			}
		})
	}()

	lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour}).
		SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"}).
		SetMethods([]string{"GET", "POST", "DELETE", "UPDATE"}).
		SetTokenBucketExpirationTTL(time.Hour)

	lib.InitializeRoutes(lmt)

	fmt.Printf("Starting server at port %s\n", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func initDB() error {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	err := mgm.SetDefaultConfig(nil, "screen_recorder", clientOptions)
	if err != nil {
		return err
	}

	client, err := mgm.NewClient(clientOptions)
	if err != nil {
		return err
	}

	return client.Ping(context.TODO(), nil)
}
