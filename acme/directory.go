package acme

import (
	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func GetDirectory(client MQTT.Client) {

	//build json
	json, _ := json.Marshal(&Directory{
		Path: "/acme/acme/directory",
	})

	//if err != nil {
	//	fmt.Println("could not marshal json: %s\n", err)
	//	return nil
	//}

	// publish getDirectory request to servers endpoint. Receive response in handler
	PublishMQTTMessage(client, json, "/acme/server")
}
