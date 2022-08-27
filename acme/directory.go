package acme

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func GetDirectory(client MQTT.Client) []byte {
	fmt.Printf("GetDirectory\n")

	//build json
	json, _ := json.Marshal(&DirectoryReq{
		Path: "/acme/acme/directory",
	})

	//if err != nil {
	//	fmt.Println("could not marshal json: %s\n", err)
	//	return nil
	//}

	// publish getDirectory request to servers endpoint. Receive response in handler
	return json
}

//parses the MQTT directory response from the acme server into a struct object
func HandleDirectoryResponse(msg MQTT.Message) *DirectoryResp {
	fmt.Printf("HandleDirectoryResponse\n")
	var dir DirectoryResp
	err := json.Unmarshal(msg.Payload(), &dir)
	if err != nil {
		fmt.Printf("couldn't parse msg into DirectoryResp %s\n", err.Error())
	}

	fmt.Printf("%+v\n", dir)
	fmt.Printf(string(msg.Payload()) + "\n")
	return &dir
}
