package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

type MessagingClientManager interface {
	SendDocument(document, recipientID, caption string, link bool) (map[string]interface{}, error)
	SendMessageText(message, recipientID string) (map[string]interface{}, error)
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

func (c *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (c *Controller) parsingMessage(message []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Printf("Error parsing message: %v, data:%v", err, data)
		return err
	}
	changedField := messengerChangedField(data)

	if changedField == "messages" {
		newMessage, _ := messengerIsMessage(data)
		if newMessage {
			mobile, _ := getMobile(data)
			name, _ := getName(data)
			messageType, _ := getMessageType(data)
			log.Printf("New Message; sender:%s name:%s type:%s", mobile, name, messageType)
			//https://whatsapp-businessapi.herokuapp.com/api/v1/4917635163191/document
			//"https://whatsapp-businessapi.herokuapp.com/api/v1/4917635163191/document/"
			url := fmt.Sprintf("https://whatsapp-businessapi.herokuapp.com/api/v1/%s/document", mobile)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				msg := fmt.Sprintf("Hormetli %s.Analiz neticeleriniz hazir degildir", name)
				_, err = c.messagingClientManager.SendMessageText(msg, mobile)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			} else {
				caption := fmt.Sprintf("Hormetli %s. Analiz neticeleriniz hazirdir", name)
				_, err = c.messagingClientManager.SendDocument(url, mobile, caption, true)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			}

		} else {
			delivery, _ := messengerGetDelivery(data)
			if delivery != nil {
				log.Printf("Message : %v", delivery)
			} else {
				log.Println("No new message")
			}
		}
	}

	return nil
}

func (c *Controller) ReceiveMessage(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(string(bytes))
	err = c.parsingMessage(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")

}

func (c *Controller) VerifyToken(w http.ResponseWriter, r *http.Request) {
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

// RequestData represents the JSON data structure
type RequestData struct {
	Number   string `json:"number"`
	Document string `json:"document"`
}

func (c *Controller) UploadDocument(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON data
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the base64-encoded PDF
	pdfData, err := base64.StdEncoding.DecodeString(requestData.Document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the PDF to disk with the filename as the number
	filename := fmt.Sprintf("%s.pdf", requestData.Number)
	err = ioutil.WriteFile(filename, pdfData, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
}

func messengerChangedField(data map[string]interface{}) string {
	if _, ok := data["entry"]; !ok {
		return ""
	}

	entry := data["entry"].([]interface{})
	changes := entry[0].(map[string]interface{})["changes"].([]interface{})
	field := changes[0].(map[string]interface{})["field"].(string)
	return field
}

func messengerGetDelivery(data map[string]interface{}) (interface{}, error) {
	preprocessedData, err := messengerPreprocess(data)
	if err != nil {
		return nil, err
	}

	if statuses, ok := preprocessedData["statuses"].([]interface{}); ok {
		if len(statuses) > 0 {
			status := statuses[0].(map[string]interface{})["status"]
			return status, nil
		}
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

func (c *Controller) GetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["number"]

	// Create a file path based on the user ID
	filePath := fmt.Sprintf("%s.pdf", number) // Replace with the actual path to your documents folder

	log.Println(filePath)
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	http.ServeFile(w, r, filePath)
}
