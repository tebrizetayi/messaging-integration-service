package api

import "testing"

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
									"wa_id": "4917635163191"
								}
							],
							"messages": [
								{
									"from": "4917635163191",
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
