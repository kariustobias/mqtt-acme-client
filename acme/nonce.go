package acme

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func HandleNewNonceRequest(client MQTT.Client, path string) []byte {
	fmt.Printf("HandleNewNonceRequest\n")
	fmt.Printf(path + "\n")
	//build json
	json, _ := json.Marshal(&NewNonceReq{
		Path: path,
	})

	// publish getDirectory request to servers endpoint. Receive response in handler
	return json
}
