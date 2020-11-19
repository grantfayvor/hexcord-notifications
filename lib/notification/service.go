package notification

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

//Notification concrete struct
type Notification struct {
	Message          map[string]interface{} `json:"message" bson:"message"`
	Recipients       []string               `json:"recipients" bson:"recipients"`
	mgm.DefaultModel `bson:",inline"`
}

//GetRecipients implemented method on Notification object
func (n *Notification) GetRecipients() []string {
	return n.Recipients
}

//GetMessage implemented method on Notification object
func (n *Notification) GetMessage() map[string]string {
	result := make(map[string]string)

	for key, val := range n.Message {
		oid, ok := val.(primitive.ObjectID)
		if ok {
			result[key] = oid.Hex()
		} else {
			result[key] = val.(string)
		}
	}
	return result
}

//UnmarshalJSON method used to customize the behaviour of json.unmarshal
func (n *Notification) UnmarshalJSON(data []byte) error {
	v := make(map[string]interface{})

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v["message"] == nil {
		return errors.New("A malformed notification request was sent")
	}

	if v["recipients"] == nil {
		v["recipients"] = make([]interface{}, 0)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(recipients []interface{}) {
		n.Recipients = []string{}

		for _, val := range recipients {
			n.Recipients = append(n.Recipients, val.(string))
		}
		wg.Done()
	}(v["recipients"].([]interface{}))

	go func(message map[string]interface{}) {
		n.Message = make(map[string]interface{})

		for key, val := range message {
			hex, err := primitive.ObjectIDFromHex(val.(string))
			if err != nil {
				n.Message[key] = val.(string)
			} else {
				n.Message[key] = hex
			}
		}
		wg.Done()
	}(v["message"].(map[string]interface{}))

	wg.Wait()

	return nil
}

//SaveNotification method used to store notification info in the database
func SaveNotification(notification *Notification) error {
	return mgm.Coll(notification).Create(notification)
}

//GetUserNotifications method for fetching user notifications
func GetUserNotifications(userID primitive.ObjectID) (result []*Notification, err error) {
	err = mgm.Coll(&Notification{}).SimpleFind(&result, bson.M{"message.recipient": userID})
	return
}
