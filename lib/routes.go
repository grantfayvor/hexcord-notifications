package lib

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/grantfayvor/hexcord-notifications/lib/notification"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func verifyAuth(w http.ResponseWriter, r *http.Request) (user map[string]interface{}, err error) {
	authorization := r.Header.Get("Authorization")

	request, err := http.NewRequest("GET", os.Getenv("HEXCORD_MASTER_ENDPOINT")+"/oauth/verify_token", nil)
	if err != nil {
		return
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", authorization)

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&user)
	return
}

func getUserNotifications(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	user, err := verifyAuth(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(user["_id"].(string))
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	notifications, err := notification.GetUserNotifications(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifications)
}

//InitializeRoutes exported function to intiailize routes
func InitializeRoutes() {
	http.HandleFunc("/user/notifications", getUserNotifications)
}
