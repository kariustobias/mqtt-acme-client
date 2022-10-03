package acme

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func HandleNewOrderRequest(client MQTT.Client, path string, jws *JWS, serialNumber string) []byte {
	fmt.Printf("HandleNewOrderRequest\n")

	var identifiers []Identifier
	identifiers = append(identifiers, Identifier{
		Type:  "permanent-identifier",
		Value: serialNumber,
	})

	payload := NewOrderReq{Identifiers: identifiers}

	payloadBytes, _ := json.Marshal(payload)

	signedContent, _ := jws.SignContent(path, payloadBytes)

	return []byte(signedContent.FullSerialize())
}

func HandleNewOrderResponse(msg MQTT.Message) *NewOrderResp {
	fmt.Printf("HandleNewOrderResponse\n")
	core.nonceManager.GetAndPushFromResponse(msg)

	var order NewOrderResp
	err := json.Unmarshal(msg.Payload(), &order)
	if err != nil {
		fmt.Printf("couldn't parse msg into NewOrderResp %s\n", err.Error())
	}
	return &order
}
