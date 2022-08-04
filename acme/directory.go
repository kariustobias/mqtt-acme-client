package acme

func GetDirectory() {
	// publish getDirectory request to servers endpoint. Receive response in handler
	PublishMQTTMessage(nil, "/acme/server")
}
