package whatsapp

import (
	"fmt"
	"log"
	"testing"
)

func TestHelloMessage_Template_Success(t *testing.T) {
	// Arrange
	client, err := NewClient("552041023667800", "e6de5aff86bed1577c681e73edf30f7e", "https://graph.facebook.com/oauth/access_token", "https://graph.facebook.com/v16.0/")
	if err != nil {
		t.Logf("error creating client: %v", err)
	}

	to := "994552178732"
	templateName := "hello_world"
	languageCode := "en_US"

	_, err = client.SendMessage(to, templateName, languageCode)
	if err != nil {
		t.Logf("error sending message: %v", err)
	}
	/*assert.Equal(t, resp.Contacts[0].WaID, to)
	assert.Equal(t, resp.Contacts[0].Input, to)
	assert.Equal(t, resp.MessagingProduct, "whatsapp")
	assert.Equal(t, resp.Messages[0].ID, "0")*/
}

func TestSendMessageTexSuccess(t *testing.T) {
	// Arrange
	client, err := NewClient("994552178732", "e6de5aff86bed1577c681e73edf30f7e", "https://graph.facebook.com/oauth/access_token", "https://graph.facebook.com/v16.0/")
	if err != nil {
		t.Logf("error creating client: %v", err)
	}

	message := "Huseyn necedir?"
	recipientID := "4917635163191"
	recipientType := "individual"

	response, err := client.SendMessageText(message, recipientID, recipientType)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(response)
}

func TestSendMessageDocument_Success(t *testing.T) {
	// Arrange
	client, err := NewClient("552041023667800", "e6de5aff86bed1577c681e73edf30f7e", "https://graph.facebook.com/oauth/access_token", "https://graph.facebook.com/v16.0/")
	if err != nil {
		t.Logf("error creating client: %v", err)
	}

	recipientID := "4917635163191"
	document := "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
	caption := ""
	link := true

	_, err = client.SendDocument(document, recipientID, caption, link)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}
