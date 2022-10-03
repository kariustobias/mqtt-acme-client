package acme

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func HandleFinalizeRequest(path string, jws *JWS, csr []byte) []byte {
	fmt.Printf("HandleFinalizeRequest\n")

	payload := FinalizeReq{
		CSR: base64.RawURLEncoding.EncodeToString(csr),
	}
	payloadBytes, _ := json.Marshal(payload)

	signedContent, _ := jws.SignContent(path, payloadBytes)

	return []byte(signedContent.FullSerialize())
}

func HandleFinalizeResponse(msg MQTT.Message) *FinalizeResp {
	fmt.Printf("HandleFinalizeResponse\n")
	core.nonceManager.GetAndPushFromResponse(msg)

	var finalize FinalizeResp
	err := json.Unmarshal(msg.Payload(), &finalize)
	if err != nil {
		fmt.Printf("couldn't parse msg into FinalizeResp %s\n", err.Error())
	}
	return &finalize
}

func CreateCSR(privateKey crypto.PrivateKey, identifier string, san []string) []byte {
	template := x509.CertificateRequest{
		Subject:  pkix.Name{CommonName: identifier},
		DNSNames: san,
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		fmt.Printf("error while generating csr : %s\n", csr)
	}

	return csr
}
