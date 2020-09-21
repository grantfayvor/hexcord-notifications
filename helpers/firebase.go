package helpers

import (
	"context"
	"os"

	fb "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

//Firebase util object
type Firebase struct {
	App *fb.App
}

//INotification interface used to send notifications through message queues
type INotification interface {
	GetMessage() map[string]string
	GetRecipients() []string
}

//Notification concrete struct
type Notification struct {
	message    map[string]string
	recipients []string
}

//NewNotification constructor object used to create Notifications
func NewNotification(message map[string]string, recipients []string) Notification {
	return Notification{message, recipients}
}

//GetMessage implemented method on Notification object
func (n *Notification) GetMessage() map[string]string {
	return n.message
}

//GetRecipients implemented method on Notification object
func (n *Notification) GetRecipients() []string {
	return n.recipients
}

//InitApp method for firebase
func (f *Firebase) InitApp() (*Firebase, error) {
	app, err := fb.NewApp(context.Background(), &fb.Config{
		ProjectID:     os.Getenv("FIREBASE_PROJECT_ID"),
		DatabaseURL:   os.Getenv("FIREBASE_DATABASE_URL"),
		StorageBucket: os.Getenv("FIREBASE_STORAGE_BUCKET"),
	}, option.WithAPIKey(os.Getenv("FIREBASE_API_KEY")))

	if err != nil {
		return nil, err
	}

	f.App = app
	return f, nil
}

//PushNotification method used to send messages to recipients
func (f *Firebase) PushNotification(notification INotification, recipient string) (string, error) {
	client, err := f.App.Messaging(context.Background())
	if err != nil {
		return "", err
	}

	msg := &messaging.Message{Data: notification.GetMessage(), Token: recipient}
	return client.Send(context.Background(), msg)
}
