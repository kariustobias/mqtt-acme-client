package acme

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	flag          bool = false
	core          *Core
	orderResponse *NewOrderResp
	serialNumber  = "1234567"
	useSimulator  = false
)

type Core struct {
	directoryResponse *DirectoryResp
	nonceManager      *Manager
	jws               *JWS
	publicKey         rsa.PublicKey
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

	c := &Core{nonceManager: nonceManager, jws: jws, directoryResponse: directoryResponse, publicKey: privateKey.PublicKey}

	return c
}

func connectToBroker(broker string) (MQTT.Client, error) {
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID("acme-client")
	opts.SetDefaultPublishHandler(f)
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil,
			fmt.Errorf("failed create new MQTT client: %w", token.Error())
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
	var requestPath Path

	// get json out of msg
	json.Unmarshal(msg.Payload(), &path)
	var json []byte
	fmt.Printf("path received: %s\n", path.Path)

	switch {
	case path.Path == "/acme/acme/directory":
		directoryResponse := HandleDirectoryResponse(msg)
		ValidateCertificate(directoryResponse.Cert)
		core = CreateCore(directoryResponse.NewNonce, directoryResponse)
		json = HandleNewNonceRequest(client, directoryResponse.NewNonce)
	case path.Path == "/acme/acme/new-nonce":
		core.nonceManager.GetAndPushFromResponse(msg)
		json = HandleNewAccountRequest(client, core.directoryResponse.NewAccount, core.jws)
		json = appendPath(requestPath, "/acme/acme/new-account", json)
	case path.Path == "/acme/acme/new-account":
		HandleNewAccountResponse(core, msg)
		json = HandleNewOrderRequest(client, core.directoryResponse.NewOrder, core.jws, serialNumber)
		json = appendPath(requestPath, "/acme/acme/new-order", json)
	case path.Path == "/acme/acme/new-order":
		orderResponse = HandleNewOrderResponse(msg)
		json = HandleNewAuthorizationRequest(orderResponse.Authorizations[0], core.jws)
		json = appendPath(requestPath, orderResponse.Authorizations[0], json)
	case strings.Contains(path.Path, "/acme/authz"):
		authzResponse := HandleNewAuthorizationResponse(msg)
		//need path in HandleGetChallengeRequest as well
		json = HandleGetChallengeRequest(authzResponse.Challenges[0].URL, core.jws, core.publicKey)
		//path is challenge id in "url" field from authorization response
		json = appendPath(requestPath, authzResponse.Challenges[0].URL, json)
	case strings.Contains(path.Path, "/acme/challenge"):
		HandleGetAndValidateChallengeResponse(msg)
		csr := CreateCSR(core.jws.privKey, serialNumber, []string{serialNumber})
		json = HandleFinalizeRequest(orderResponse.Finalize, core.jws, csr)
		json = appendPath(requestPath, orderResponse.Finalize, json)
	case strings.Contains(path.Path, "finalize"):
		finalizeResponse := HandleFinalizeResponse(msg)
		json = HandleCertRequest(finalizeResponse.Certificate, core.jws)
		json = appendPath(requestPath, finalizeResponse.Certificate, json)
	default:
		fmt.Printf("undefined path value\n")
		fmt.Printf(path.Path + "\n")
		fmt.Println(string(msg.Payload()) + "\n")
		json = nil
	}

	json = EncryptData(json, core.directoryResponse.Cert)
	//PublishJson
	PublishMQTTMessage(client, json, "/acme/server")

	// store json as type maybe

	//set flag to true, in the end
	//flag = true
}

// called on response to GetDirectory request from server
// sends out a new-nonce request

func appendPath(requestPath Path, path string, json []byte) []byte {
	requestPath.Path = path
	reqPath := fmt.Sprintf(",\"path\":\"%s\"}", requestPath.Path)
	return append(json[:len(json)-1], reqPath...)
}

func HandleCertRequest(path string, jws *JWS) []byte {
	fmt.Printf("HandleCertRequest\n")
	time.Sleep(2 * time.Second)

	var payloadBytes = []byte{}

	signedContent, _ := jws.SignContent(path, payloadBytes)

	return []byte(signedContent.FullSerialize())
}
