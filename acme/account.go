package acme

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func HandleNewAccountRequest(client MQTT.Client, path string, jws *JWS) []byte {
	fmt.Printf("HandleNewAccountRequest\n")

	// create Account object
	payload := NewAccountReq{
		TermsOfServiceAgreed: true,
		Contact:              []string{"mailto:tobias.karius@yahoo.de"},
	}
	payloadBytes, _ := json.Marshal(payload)

	signedContent, _ := jws.SignContent(path, payloadBytes)

	return []byte(signedContent.FullSerialize())
}

func HandleNewAccountResponse(core *Core, msg MQTT.Message) {
	// push nonce first
	core.nonceManager.GetAndPushFromResponse(msg)
	// initialize kid. Set it to "orders" field

	// 1. parse msg in NewAccountResp type struct
	var accountResponse NewAccountResp
	err := json.Unmarshal(msg.Payload(), &accountResponse)
	if err != nil {
		fmt.Printf("couldn't parse msg into NewAccountResp %s\n", err.Error())
	}
	// 2. set kid
	core.jws.kid = accountResponse.Orders
}
