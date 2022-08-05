package acme

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//takes a marshaled json as input parameter and publishes the json as string on given endpoint
func PublishMQTTMessage(client MQTT.Client, json []byte, endpoint string) {
	//string(json)
	token := client.Publish(endpoint, 0, false, string(json))
	token.Wait()
}
