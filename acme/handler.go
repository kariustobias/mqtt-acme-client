package acme

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

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

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {

}
