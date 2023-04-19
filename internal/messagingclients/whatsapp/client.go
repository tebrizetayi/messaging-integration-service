package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	phoneID = "106189092448679"

	MessagingProduct    = "whatsapp"
	MessageTypeTemplate = "template"
	MessageTypeText     = "text"
	MessageTypeDocument = "document"

	RequestTypeIndividual = "individual"
)

type Client struct {
	ClientID          string
	AccessToken       string
	accessTokenURL    string
	sendingMessageURL string
	BearerToken       string
}

func NewClient(clientID, accessToken, accessTokenURL, sendingMessageURL, bearerToken string) Client {
	return Client{
		ClientID:          clientID,
		AccessToken:       accessToken,
		accessTokenURL:    accessTokenURL,
		sendingMessageURL: sendingMessageURL,
		BearerToken:       bearerToken,
	}
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

func (c *Client) GetUrl() string {
	return fmt.Sprintf("%s%s/messages", c.sendingMessageURL, phoneID)
}

func (c *Client) SendMessage(to, templateName, languageCode string) (SendMessageResponse, error) {
	url := c.GetUrl()
	payload := SendMessagePayload{
		MessagingProduct: MessagingProduct,
		To:               to,
		Type:             MessageTypeTemplate,
		Template: Template{
			Name:     templateName,
			Language: TemplateLanguage{Code: languageCode},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return SendMessageResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return SendMessageResponse{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.BearerToken))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SendMessageResponse{}, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SendMessageResponse{}, err
	}

	var sendMessageResponse SendMessageResponse
	err = json.Unmarshal(data, &sendMessageResponse)
	if err != nil {
		return sendMessageResponse, err
	}

	return sendMessageResponse, nil
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

func (c *Client) SendMessageText(message, recipientID string) (map[string]interface{}, error) {
	url := c.GetUrl()

	data := SendMessageText{
		MessagingProduct: MessagingProduct,
		RecipientType:    RequestTypeIndividual,
		To:               recipientID,
		Type:             MessageTypeText,
		Text:             Text{PreviewURL: false, Body: message},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.BearerToken))
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
	url := c.GetUrl()
	data := SendDocumentRequest{
		MessagingProduct: MessagingProduct,
		To:               recipientID,
		Type:             MessageTypeDocument,
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
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.BearerToken))
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

	return result, nil
}
