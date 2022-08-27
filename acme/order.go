package acme

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func HandleNewOrderRequest(client MQTT.Client, path string, jws *JWS) []byte {
	fmt.Printf("HandleNewOrderRequest\n")

	var identifiers []Identifier
	identifiers = append(identifiers, Identifier{Type: "dns", Value: "test.test.com"})

	payload := NewOrderReq{Identifiers: identifiers}

	payloadBytes, _ := json.Marshal(payload)

	signedContent, _ := jws.SignContent(path, payloadBytes)

	return []byte(signedContent.FullSerialize())
}
