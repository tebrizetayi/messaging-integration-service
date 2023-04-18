package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	phoneID = "106189092448679"
)

type Client struct {
	ClientID          string
	AccessToken       string
	accessTokenURL    string //https://graph.facebook.com/oauth/access_token?client_id=%s&client_secret=%s&grant_type=client_credentials"
	sendingMessageURL string //https://graph.facebook.com/v16.0/

}

func NewClient(clientID string, accessToken string, accessTokenURL string, sendingMessageURL string) (Client, error) {
	client := Client{
		ClientID:          clientID,
		AccessToken:       accessToken,
		accessTokenURL:    accessTokenURL,
		sendingMessageURL: sendingMessageURL,
	}

	return client, nil
}

func (c *Client) getBearerToken() string {
	//return ""
	return os.Getenv("WHATSAPP_ACCESS_TOKEN")
}

type TemplateLanguage struct {
	Code string `json:"code"`
}

type Template struct {
	Name     string           `json:"name"`
	Language TemplateLanguage `json:"language"`
}

type SendMessagePayload struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

type SendMessageResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

// url = "https://graph.facebook.com/v16.0/106189092448679/messages"
func (c *Client) SendMessage(to string, templateName string, languageCode string) (*SendMessageResponse, error) {
	url := fmt.Sprintf("%s%s/messages", c.sendingMessageURL, phoneID)

	payload := SendMessagePayload{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: Template{
			Name:     templateName,
			Language: TemplateLanguage{Code: languageCode},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jsonPayload))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.getBearerToken())
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := []byte{}
	_, err = resp.Body.Read(data)
	if err != nil {
		return nil, err
	}

	var sendMessageResponse SendMessageResponse
	err = json.Unmarshal(data, &sendMessageResponse)
	if err != nil {
		return nil, err
	}

	return &sendMessageResponse, nil
}

//curl -i -X POST \
//https://graph.facebook.com/v16.0/105954558954427/messages \
// -H 'Authorization: Bearer EAAFl...' \
// -H 'Content-Type: application/json' \
// -d '{ "messaging_product": "whatsapp", "to": "15555555555", "type": "template", "template": { "name": "hello_world", "language": { "code": "en_US" } } }'

func (c *Client) SendCustomMessage(to string, templateName string, languageCode string) (*SendMessageResponse, error) {
	url := fmt.Sprintf("%s%s/messages", c.sendingMessageURL, phoneID)

	payload := SendMessagePayload{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: Template{
			Name:     templateName,
			Language: TemplateLanguage{Code: languageCode},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jsonPayload))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.getBearerToken())
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := []byte{}
	_, err = resp.Body.Read(data)
	if err != nil {
		return nil, err
	}

	var sendMessageResponse SendMessageResponse
	err = json.Unmarshal(data, &sendMessageResponse)
	if err != nil {
		return nil, err
	}

	return &sendMessageResponse, nil
}

type Text struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type SendMessageText struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             Text   `json:"text"`
}

func (c *Client) SendMessageText(message, recipientID, recipientType string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s/messages", c.sendingMessageURL, phoneID)
	data := SendMessageText{
		MessagingProduct: "whatsapp",
		RecipientType:    recipientType,
		To:               recipientID,
		Type:             "text",
		Text:             Text{PreviewURL: false, Body: message},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Printf("Sending message to %s", recipientID)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.getBearerToken())
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		log.Printf("Message sent to %s", recipientID)
	} else {
		log.Printf("Message not sent to %s", recipientID)
		log.Printf("Status code: %d", resp.StatusCode)
		log.Printf("Response: %v", result)
	}

	return result, nil
}

type Document struct {
	Link    string `json:"link,omitempty"`
	ID      string `json:"id,omitempty"`
	Caption string `json:"caption,omitempty"`
}

type SendDocumentRequest struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Document         Document `json:"document"`
}

func (c *Client) SendDocument(document, recipientID, caption string, link bool) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s/messages", c.sendingMessageURL, phoneID)
	data := SendDocumentRequest{
		MessagingProduct: "whatsapp",
		To:               recipientID,
		Type:             "document",
	}

	if link {
		data.Document = Document{Link: document, Caption: caption}
	} else {
		data.Document = Document{ID: document, Caption: caption}
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Printf("Sending document to %s", recipientID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.getBearerToken())
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	log.Println(string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		log.Printf("Document sent to %s", recipientID)
	} else {
		log.Printf("Document not sent to %s", recipientID)
		log.Printf("Status code: %d", resp.StatusCode)
		log.Printf("Response: %v", result)
	}

	return result, nil
}
