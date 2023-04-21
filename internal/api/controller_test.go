package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestParseMessage_Success(t *testing.T) {
	c := NewController(nil)
	data := []byte(`{
		"object": "whatsapp_business_account",
		"entry": [
			{
				"id": "102140959526615",
				"changes": [
					{
						"value": {
							"messaging_product": "whatsapp",
							"metadata": {
								"display_phone_number": "15550909792",
								"phone_number_id": "106189092448679"
							},
							"contacts": [
								{
									"profile": {
										"name": "T.A"
									},
									"wa_id": "994503981865"
								}
							],
							"messages": [
								{
									"from": "994503981865",
									"id": "wamid.HBgNNDkxNzYzNTE2MzE5MRUCABIYFjNFQjA4QTQyQTIxOUZCMDI2QTFDRTYA",
									"timestamp": "1681899808",
									"text": {
										"body": "test"
									},
									"type": "text"
								}
							]
						},
						"field": "messages"
					}
				]
			}
		]
	}`)
	err := c.parsingMessage(data)
	if err != nil {
		t.Errorf("error parsing message: %v", err)
	}

}

func TestParseMessage_Example(t *testing.T) {
	c := NewController(nil)
	data := []byte(`{"messaging_product":"whatsapp","contacts":[{"input":"4917635163191","wa_id":"4917635163191"}],"messages":[{"id":"wamid.HBgNNDkxNzYzNTE2MzE5MRUCABEYEjhDQzE0MUI5M0VBQTU4MzVBRQA="}]}`)
	err := c.parsingMessage(data)
	if err != nil {
		t.Errorf("error parsing message: %v", err)
	}
}

func TestParseMessage_Example2(t *testing.T) {
	c := NewController(nil)
	data := []byte(`{"object":"whatsapp_business_account","entry":[{"id":"102140959526615","changes":[{"value":{"messaging_product":"whatsapp","metadata":{"display_phone_number":"15550909792","phone_number_id":"106189092448679"},"statuses":[{"id":"wamid.HBgNNDkxNzYzNTE2MzE5MRUCABEYEjY3OENERDM4RTY2Mzc3RkE4MgA=","status":"delivered","timestamp":"1681903556","recipient_id":"4917635163191","conversation":{"id":"6f1a08afbb622fbac2235469f724a890","origin":{"type":"user_initiated"}},"pricing":{"billable":true,"pricing_model":"CBP","category":"user_initiated"}}]},"field":"messages"}]}]}`)
	err := c.parsingMessage(data)
	if err != nil {
		t.Errorf("error parsing message: %v", err)
	}
}

func TestUploadDocument_Success(t *testing.T) {
	pdfPath := "sample.pdf" // Replace with the path to your sample PDF file
	pdfData, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		t.Fatalf("Unable to read sample PDF file: %v", err)
	}

	encodedPDF := base64.StdEncoding.EncodeToString(pdfData)

	// Create a JSON object with the number and encoded PDF
	requestData := RequestData{
		Number:   "testnumber",
		Document: encodedPDF,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Unable to marshal JSON data: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest("POST", "/api/v1/upload-document", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	c := NewController(nil)
	api := NewAPI(c)
	// Call the uploadHandler with the test request and response recorder
	api.ServeHTTP(rr, req)
	// Check if the status code is http.StatusOK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the file was saved correctly
	if _, err := os.Stat(fmt.Sprintf("%s.pdf", requestData.Number)); os.IsNotExist(err) {
		t.Errorf("File was not saved: %v", err)
	} else {
		// Clean up the saved file
		os.Remove(fmt.Sprintf("%s.pdf", requestData.Number))
	}
}

func TestGetDocument_Success(t *testing.T) {
	pdfPath := "sample.pdf" // Replace with the path to your sample PDF file
	pdfData, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		t.Fatalf("Unable to read sample PDF file: %v", err)
	}

	encodedPDF := base64.StdEncoding.EncodeToString(pdfData)

	// Create a JSON object with the number and encoded PDF
	requestData := RequestData{
		Number:   "testnumber",
		Document: encodedPDF,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Unable to marshal JSON data: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/%s/document", "4917635163191"), strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	c := NewController(nil)
	api := NewAPI(c)
	api.ServeHTTP(rr, req)
	// Call the uploadHandler with the test request and response recorder

	// Check if the status code is http.StatusOK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
