package acme

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/go-attestation/attest"
	"github.com/google/go-tpm-tools/simulator"
	"golang.org/x/crypto/acme"
)

func HandleGetChallengeRequest(path string, jws *JWS, publicKey rsa.PublicKey) []byte {
	fmt.Printf("HandleNewChallengeRequest\n")

	//DeviceAttest01Challenge(publicKey)

	var payloadBytes = []byte{}

	signedContent, _ := jws.SignContent(path, payloadBytes)

	return []byte(signedContent.FullSerialize())
}

func HandleGetAndValidateChallengeResponse(msg MQTT.Message) *NewChallengeResp {
	fmt.Printf("HandleNewAuthorizationResponse\n")
	core.nonceManager.GetAndPushFromResponse(msg)

	var chall NewChallengeResp
	err := json.Unmarshal(msg.Payload(), &chall)
	if err != nil {
		fmt.Printf("couldn't parse msg into NewAuthorizationResp %s\n", err.Error())
	}
	if chall.Status != "valid" {
		fmt.Printf("unexpected challenge status. Expected : valid, Got : %s\n", chall.Status)
	}
	return &chall
}

type simulatorChannel struct {
	io.ReadWriteCloser
}

func (simulatorChannel) MeasurementLog() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func DeviceAttest01Challenge(useSimulator bool, publicKey rsa.PublicKey, token string) error {
	// Generate the certificate key, include the ACME key authorization in the
	// the TPM certification data.
	tpm, ak, akCert, err := tpmInit(useSimulator)
	if err != nil {
		return err
	}
	data, err := keyAuthDigest(publicKey, token)
	if err != nil {
		log.Fatal(err)
	}
	config := &attest.KeyConfig{
		Algorithm:      attest.ECDSA,
		Size:           256,
		QualifyingData: data,
	}
	certKey, err := tpm.NewKey(ak, config)
	if err != nil {
		log.Fatal(err)
	}
	// Generate the WebAuthn attestation statement.
	payload, err := attestationStatement(certKey, akCert)
	if err != nil {
		log.Fatal(err)
	}
	req := struct {
		AttStmt []byte `json:"attStmt"`
	}{
		payload,
	}
	// Fufill the ACME challenge using the attestation statement.
	chalResp, err := client.AcceptWithPayload(ctx, chal, req)
	if err != nil {
		log.Fatal(err)
	}
	if chalResp.Error != nil {
		log.Fatal(chalResp.Error)
	}
}

func tpmInit(useSimulator bool) (*attest.TPM, *attest.AK, []byte, error) {
	config := &attest.OpenConfig{}
	if useSimulator {
		sim, err := simulator.Get()

		if err != nil {
			return nil, nil, nil, err
		}
		config.CommandChannel = simulatorChannel{sim}
	}
	tpm, err := attest.OpenTPM(config)
	if err != nil {
		return nil, nil, nil, err
	}
	ak, err := tpm.NewAK(nil)
	if err != nil {
		return nil, nil, nil, err
	}
	akCert, err := akCert(ak)
	if err != nil {
		return nil, nil, nil, err
	}
	return tpm, ak, akCert, nil
}

func attestationStatement(key *attest.Key, akCert []byte) ([]byte, error) {
	params := key.CertificationParameters()

	obj := &AttestationObject{
		Format: "tpm",
		AttStatement: map[string]interface{}{
			"ver":      "2.0",
			"alg":      int64(-257), // AlgRS256
			"x5c":      []interface{}{akCert},
			"sig":      params.CreateSignature,
			"certInfo": params.CreateAttestation,
			"pubArea":  params.Public,
		},
	}
	b, err := cbor.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Borrowed from:
// https://github.com/golang/crypto/blob/master/acme/acme.go#L748
func keyAuthDigest(pub crypto.PublicKey, token string) ([]byte, error) {
	th, err := acme.JWKThumbprint(pub)
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256([]byte(fmt.Sprintf("%s.%s", token, th)))
	return digest[:], err
}

func akCert(ak *attest.AK) ([]byte, error) {
	akRootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	akRootTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
	}
	permID := x509ext.PermanentIdentifier{
		IdentifierValue: serialNumber,
		Assigner:        asn1.ObjectIdentifier{0, 1, 2, 3, 4},
	}
	san := &x509ext.SubjectAltName{
		PermanentIdentifiers: []x509ext.PermanentIdentifier{
			permID,
		},
	}
	ext, err := x509ext.MarshalSubjectAltName(san)
	if err != nil {
		return nil, err
	}
	akTemplate := &x509.Certificate{
		SerialNumber:    big.NewInt(2),
		ExtraExtensions: []pkix.Extension{ext},
	}
	akPub, err := attest.ParseAKPublic(attest.TPMVersion20, ak.AttestationParameters().Public)
	if err != nil {
		return nil, err
	}
	akCert, err := x509.CreateCertificate(rand.Reader, akTemplate, akRootTemplate, akPub.Public, akRootKey)
	if err != nil {
		return nil, err
	}
	return akCert, nil
}
