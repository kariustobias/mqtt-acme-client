package acme

import (
	"encoding/json"
	"fmt"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var flag bool = false
var core *Core

type Core struct {
	directoryResponse *DirectoryResp
	nonceManager      *Manager
	jws               *JWS
}

func Initialize(broker string, endpoint string) (MQTT.Client, error) {
	client, err := connectToBroker(broker)
	if err != nil {
		return nil, err
	}
	//subscribe
	err = subscribe(endpoint, client)
	if err != nil {
		return nil, err
	}

	//publish on GetDirectory
	return client, nil
}

//called after Directory Object is initialized
func CreateCore(newNonceUrl string, directoryResponse *DirectoryResp) *Core {
	nonceManager := NewManager(newNonceUrl)

	privateKey := GenerateKeyPair()

	//initialize with empty kid. kid will be initialized on new-account response
	jws := NewJWS(privateKey, "", nonceManager)

	c := &Core{nonceManager: nonceManager, jws: jws, directoryResponse: directoryResponse}

	return c
}

func connectToBroker(broker string) (MQTT.Client, error) {
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID("acme-client")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed create new MQTT client: %w", token.Error())
	}
	return c, nil
}

func subscribe(endpoint string, client MQTT.Client) error {
	if token := client.Subscribe(endpoint, 0, nil); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to MQTT endpoint: %w", token.Error())
	}
	return nil
}

func GetFlag() bool {
	return flag
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {

	var path Path

	// get json out of msg
	json.Unmarshal(msg.Payload(), &path)
	var json []byte
	fmt.Printf("path received: %s\n", path.Path)

	switch sub := strings.ReplaceAll(path.Path, "/acme/acme/", ""); sub {
	case "directory":
		directoryResponse := HandleDirectoryResponse(msg)
		core = CreateCore(directoryResponse.NewNonce, directoryResponse)
		json = HandleNewNonceRequest(client, directoryResponse.NewNonce)
	case "new-nonce":
		core.nonceManager.GetAndPushFromResponse(msg)
		json = HandleNewAccountRequest(client, core.directoryResponse.NewAccount, core.jws)
	case "new-account":
		HandleNewAccountResponse(core, msg)
		json = HandleNewOrderRequest(client, core.directoryResponse.NewOrder, core.jws)
	case "new-order":

	default:
		fmt.Printf("undefined path value\n")
		fmt.Printf(sub + "\n")
		fmt.Println(string(msg.Payload()) + "\n")
		json = nil
	}

	//PublishJson
	PublishMQTTMessage(client, json, "/acme/server")

	// store json as type maybe

	//set flag to true, in the end
	//flag = true
}

// called on response to GetDirectory request from server
// sends out a new-nonce request
