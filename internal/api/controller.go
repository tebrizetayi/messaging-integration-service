package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type MessagingClientManager interface {
}

// Controller is the API controller
type Controller struct {
	messagingClientManager MessagingClientManager
}

func NewController(mc MessagingClientManager) Controller {
	return Controller{
		messagingClientManager: mc,
	}
}

func (c Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (c Controller) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received webhook data: %v", data)
	changedField := messengerChangedField(data)

	if changedField == "messages" {
		newMessage, _ := messengerIsMessage(data)
		if newMessage {
			mobile, _ := getMobile(data)
			name, _ := getName(data)
			messageType, _ := getMessageType(data)
			log.Printf("New Message; sender:%s name:%s type:%s", mobile, name, messageType)

			// Process different message types here
			// e.g. text, interactive, location, image, video, audio, document

		} else {
			delivery, _ := messengerGetDelivery(data)
			if delivery != nil {
				log.Printf("Message : %v", delivery)
			} else {
				log.Println("No new message")
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (c Controller) VerifyToken(w http.ResponseWriter, r *http.Request) {
	verifyToken := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if verifyToken == os.Getenv("VERIFY_TOKEN") {
		log.Println("Verified webhook")
		w.Header().Set("Content-Type", "text/plain")

		challengeInt, err := strconv.Atoi(challenge)
		if err != nil {
			log.Printf("Error converting challenge to integer: %v", err)
			http.Error(w, "Invalid challenge value", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%d", challengeInt)
	} else {
		log.Println("Webhook Verification failed")
		http.Error(w, "Invalid verification token", http.StatusUnauthorized)
	}
}

func messengerChangedField(data map[string]interface{}) string {
	entry := data["entry"].([]interface{})
	changes := entry[0].(map[string]interface{})["changes"].([]interface{})
	field := changes[0].(map[string]interface{})["field"].(string)
	return field
}

func messengerGetDelivery(data map[string]interface{}) (map[string]interface{}, error) {
	preprocessedData, err := messengerPreprocess(data)
	if err != nil {
		return nil, err
	}

	if statuses, ok := preprocessedData["statuses"].([]interface{}); ok {
		return statuses[0].(map[string]interface{})["status"].(map[string]interface{}), nil
	}

	return nil, nil
}

func messengerIsMessage(data map[string]interface{}) (bool, error) {
	preprocessedData, err := messengerPreprocess(data)
	if err != nil {
		return false, err
	}

	if _, ok := preprocessedData["messages"]; ok {
		return true, nil
	}

	return false, nil
}

func messengerPreprocess(data map[string]interface{}) (map[string]interface{}, error) {
	entry, ok := data["entry"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("entry not found in data")
	}

	changes, ok := entry[0].(map[string]interface{})["changes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("changes not found in entry")
	}

	value, ok := changes[0].(map[string]interface{})["value"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value not found in changes")
	}

	return value, nil
}

func getMobile(data map[string]interface{}) (string, error) {
	preprocessedData, err := messengerPreprocess(data)
	if err != nil {
		return "", err
	}

	contacts, ok := preprocessedData["contacts"].([]interface{})
	if !ok {
		return "", fmt.Errorf("contacts not found in data")
	}

	waID, ok := contacts[0].(map[string]interface{})["wa_id"].(string)
	if !ok {
		return "", fmt.Errorf("wa_id not found in contacts")
	}

	return waID, nil
}

func getName(data map[string]interface{}) (string, error) {
	preprocessedData, err := messengerPreprocess(data)
	if err != nil {
		return "", err
	}

	contacts, ok := preprocessedData["contacts"].([]interface{})
	if !ok {
		return "", fmt.Errorf("contacts not found in data")
	}

	contact := contacts[0].(map[string]interface{})

	profile, ok := contact["profile"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("profile not found in contact")
	}

	name, ok := profile["name"].(string)
	if !ok {
		return "", fmt.Errorf("name not found in profile")
	}

	return name, nil
}

func getMessageType(data map[string]interface{}) (string, error) {
	preprocessedData, err := messengerPreprocess(data)
	if err != nil {
		return "", err
	}

	messages, ok := preprocessedData["messages"].([]interface{})
	if !ok {
		return "", fmt.Errorf("messages not found in data")
	}

	message := messages[0].(map[string]interface{})

	messageType, ok := message["type"].(string)
	if !ok {
		return "", fmt.Errorf("type not found in message")
	}

	return messageType, nil
}
