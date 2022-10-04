package acme

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.step.sm/crypto/pemutil"
)

//takes a marshaled json as input parameter and publishes the json as string on given endpoint
func PublishMQTTMessage(client MQTT.Client, json []byte, endpoint string) {
	//string(json)
	token := client.Publish(endpoint, 0, false, string(json))
	token.Wait()
}

func GenerateKeyPair() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("keypair generation failed: %s", err)
	}
	return privateKey
}

func ValidateCertificate(cert []byte) (bool, error) {
	c, err := x509.ParseCertificate(cert)
	if err != nil {
		fmt.Printf("certificate could not be parsed: %s", err)
		return false, err
	}

	root := x509.NewCertPool()
	crt, err := pemutil.ReadCertificate("test_certs/root_ca.crt")
	if err != nil {
		fmt.Printf("intermediate certificate could not be parsed: %s", err)
		return false, err
	}
	root.AddCert(crt)
	opts := x509.VerifyOptions{
		Roots: root,
	}

	if _, err := c.Verify(opts); err != nil {
		fmt.Printf("failed to validate certificate: %s", err)
		return false, err
	}
	return true, nil
}

func EncryptData(json []byte, cert []byte) []byte {
	c, err := x509.ParseCertificate(cert)
	if err != nil {
		fmt.Printf("failed to parse certifiate : %s", err)
	}
	rsaPublicKey := c.PublicKey.(*rsa.PublicKey)
	cipher, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, rsaPublicKey, json, nil)
	if err != nil {
		fmt.Printf("failed to encrypt json data with servers public key certificate : %s", err)
	}
	return cipher
}
