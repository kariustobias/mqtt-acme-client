package main

import (
	"github.com/kariustobias/mqtt-acme-client/acme"
)

func main() {

	// publish GETDirectory json on /acme/server endpoint
	client, _ := acme.Initialize("tcp://localhost:1883", "/acme/client")
	// publish GetDirectory request. All further requests will be handled by the handler
	acme.GetDirectory(client)
}
