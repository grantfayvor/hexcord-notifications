package notifcation

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Notification concrete struct
type Notification struct {
	Message          map[string]string    `json:"message" bson:"message"`
	Recipients       []string             `json:"recipients" bson:"recipients"`
	RecipientIDs     []primitive.ObjectID `json:"recipientIds" bson:"recipientIds"`
	mgm.DefaultModel `bson:",inline"`
}

//NewNotification constructor object used to create Notifications
func NewNotification(message map[string]string, recipients []string) Notification {
	return Notification{Message: message, Recipients: recipients}
}

//GetMessage implemented method on Notification object
func (n *Notification) GetMessage() map[string]string {
	return n.Message
}

//GetRecipients implemented method on Notification object
func (n *Notification) GetRecipients() []string {
	return n.Recipients
}

//SaveNotification method used to store notification info in the database
func SaveNotification(notification *Notification) error {
	for i, id := range notification.Recipients {
		_id, _ := primitive.ObjectIDFromHex(id)
		notification.RecipientIDs[i] = _id
	}
	return mgm.Coll(notification).Create(notification)
}
