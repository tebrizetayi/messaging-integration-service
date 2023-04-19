package whatsapp

import (
	"fmt"
	"log"
	"testing"
)

func TestHelloMessage_Template_Success(t *testing.T) {
	// Arrange
	client := NewClient("552041023667800", "e6de5aff86bed1577c681e73edf30f7e", "https://graph.facebook.com/v16.0/", "")

	to := "994552178732"
	templateName := "hello_world"
	languageCode := "en_US"

	_, err := client.SendMessage(to, templateName, languageCode)
	if err != nil {
		t.Logf("error sending message: %v", err)
	}
}

func TestSendMessageTexSuccess(t *testing.T) {
	// Arrange
	client := NewClient("994552178732", "e6de5aff86bed1577c681e73edf30f7e", "https://graph.facebook.com/v16.0/", "")

	message := "It is a Test message.Just ignore it."
	recipientID := "4917635163191"

	response, err := client.SendMessageText(message, recipientID)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(response)
}

func TestSendMessageDocument_Success(t *testing.T) {
	// Arrange
	client := NewClient("552041023667800", "e6de5aff86bed1577c681e73edf30f7e", "https://graph.facebook.com/v16.0/", "")

	recipientID := "4917635163191"
	document := "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
	caption := ""
	link := true

	_, err := client.SendDocument(document, recipientID, caption, link)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}
