package main

import (
	"time"

	"github.com/kariustobias/mqtt-acme-client/acme"
)

func main() {

	// publish GETDirectory json on /acme/server endpoint
	client, _ := acme.Initialize("tcp://localhost:1883", "/acme/client")
	// publish GetDirectory request. All further requests will be handled by the handler
	json := acme.GetDirectory(client)
	acme.PublishMQTTMessage(client, json, "/acme/server")

	nonceManager := nonces.NewManager(doer, dir.NewNonceURL)

	jws := secure.NewJWS(privateKey, kid, nonceManager)

	for !acme.GetFlag() {
		time.Sleep(1 * time.Second)
		//fmt.Println("waiting: ", wcount)
		//wcount += 1
	}
}
