package acme

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func HandleNewAuthorizationRequest(path string, jws *JWS) []byte {
	fmt.Printf("HandleNewAuthorizationRequest\n")

	var payloadBytes = []byte{}

	signedContent, _ := jws.SignContent(path, payloadBytes)

	return []byte(signedContent.FullSerialize())
}

func HandleNewAuthorizationResponse(msg MQTT.Message) *NewAuthorizationResp {
	fmt.Printf("HandleNewAuthorizationResponse\n")
	core.nonceManager.GetAndPushFromResponse(msg)

	var authz NewAuthorizationResp
	err := json.Unmarshal(msg.Payload(), &authz)
	if err != nil {
		fmt.Printf("couldn't parse msg into NewAuthorizationResp %s\n", err.Error())
	}
	return &authz
}
