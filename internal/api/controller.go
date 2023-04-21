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
	SendDocument(from, document, recipientID, caption string, link bool) (map[string]interface{}, error)
	SendMessageText(from, message, recipientID string) (map[string]interface{}, error)
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

	md := &MessengerData{
		preprocessedData: data,
	}
	changedField := md.ChangedField()

	if changedField == "messages" {
		newMessage := md.IsMessage()
		if newMessage {
			mobile, _ := md.GetMobile()
			name, _ := md.GetName()
			messageType, _ := md.GetMessageType()
			businessNumber, _ := md.GetBusinessNumber()

			log.Printf("New Message; sender:%s name:%s type:%s", mobile, name, messageType)
			url := fmt.Sprintf("https://whatsapp-businessapi.herokuapp.com/api/v1/%s/document", mobile)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}

			_ = businessNumber

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				msg := fmt.Sprintf("Hormetli %s.Analiz neticeleriniz hazir degildir", name)
				_, err = c.messagingClientManager.SendMessageText(businessNumber, msg, mobile)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			} else {
				caption := fmt.Sprintf("Hormetli %s. Analiz neticeleriniz hazirdir", name)
				_, err = c.messagingClientManager.SendDocument(businessNumber, url, mobile, caption, true)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			}

		} else {
			delivery, _ := md.GetDelivery()
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
