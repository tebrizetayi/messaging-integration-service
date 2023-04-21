package api

import (
	"errors"
)

type MessengerData struct {
	preprocessedData map[string]interface{}
}

func (md *MessengerData) Preprocess(inputData map[string]interface{}) error {
	entryList, ok := inputData["entry"].([]interface{})
	if !ok {
		return errors.New("entry not found in data")
	}

	changeList, ok := entryList[0].(map[string]interface{})["changes"].([]interface{})
	if !ok {
		return errors.New("changes not found in entry")
	}

	valueMap, ok := changeList[0].(map[string]interface{})["value"].(map[string]interface{})
	if !ok {
		return errors.New("value not found in changes")
	}

	md.preprocessedData = valueMap
	return nil
}

func (md *MessengerData) ChangedField() string {
	if _, ok := md.preprocessedData["entry"]; !ok {
		return ""
	}

	entryList := md.preprocessedData["entry"].([]interface{})
	changeList := entryList[0].(map[string]interface{})["changes"].([]interface{})
	fieldName := changeList[0].(map[string]interface{})["field"].(string)
	return fieldName
}

func (md *MessengerData) GetDelivery() (interface{}, error) {
	if statusList, ok := md.preprocessedData["statuses"].([]interface{}); ok {
		if len(statusList) > 0 {
			deliveryStatus := statusList[0].(map[string]interface{})["status"]
			return deliveryStatus, nil
		}
	}

	return nil, nil
}

func (md *MessengerData) IsMessage() bool {
	_, ok := md.preprocessedData["messages"]
	return ok
}

func (md *MessengerData) GetMobile() (string, error) {
	contactList, ok := md.preprocessedData["contacts"].([]interface{})
	if !ok {
		return "", errors.New("contacts not found in data")
	}

	waID, ok := contactList[0].(map[string]interface{})["wa_id"].(string)
	if !ok {
		return "", errors.New("wa_id not found in contacts")
	}

	return waID, nil
}

func (md *MessengerData) GetName() (string, error) {
	contactList, ok := md.preprocessedData["contacts"].([]interface{})
	if !ok {
		return "", errors.New("contacts not found in data")
	}

	contactInfo := contactList[0].(map[string]interface{})

	profileInfo, ok := contactInfo["profile"].(map[string]interface{})
	if !ok {
		return "", errors.New("profile not found in contact")
	}

	contactName, ok := profileInfo["name"].(string)
	if !ok {
		return "", errors.New("name not found in profile")
	}

	return contactName, nil
}

func (md *MessengerData) GetMessageType() (string, error) {
	messageList, ok := md.preprocessedData["messages"].([]interface{})
	if !ok {
		return "", errors.New("messages not found in data")
	}

	messageInfo := messageList[0].(map[string]interface{})

	messageType, ok := messageInfo["type"].(string)
	if !ok {
		return "", errors.New("type not found in message")
	}

	return messageType, nil
}

func (md *MessengerData) GetBusinessNumber() (string, error) {
	metaInfo, ok := md.preprocessedData["metadata"].(map[string]interface{})
	if !ok {
		return "", errors.New("meta not found in data")
	}

	phoneNumber, ok := metaInfo["display_phone_number"].(string)
	if !ok {
		return "", errors.New("wa_id not found in contacts")
	}

	return phoneNumber, nil
}
