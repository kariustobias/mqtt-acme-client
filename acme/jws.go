package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"

	"go.step.sm/crypto/jose"
)

type JWS struct {
	privKey crypto.PrivateKey
	kid     string // Key identifier
	nonces  *Manager
}

func NewJWS(privateKey crypto.PrivateKey, kid string, nonceManager *Manager) *JWS {
	return &JWS{
		privKey: privateKey,
		nonces:  nonceManager,
		kid:     kid,
	}
}

func (j *JWS) SignContent(url string, content []byte) (*jose.JSONWebSignature, error) {
	var alg jose.SignatureAlgorithm
	switch k := j.privKey.(type) {
	case *rsa.PrivateKey:
		alg = jose.RS256
	case *ecdsa.PrivateKey:
		if k.Curve == elliptic.P256() {
			alg = jose.ES256
		} else if k.Curve == elliptic.P384() {
			alg = jose.ES384
		}
	}

	signKey := jose.SigningKey{
		Algorithm: alg,
		Key:       jose.JSONWebKey{Key: j.privKey, KeyID: j.kid},
	}

	options := jose.SignerOptions{
		NonceSource: j.nonces,
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"url": url,
		},
	}

	if j.kid == "" {
		options.EmbedJWK = true
	}

	signer, err := jose.NewSigner(signKey, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to create jose signer: %w", err)
	}

	signed, err := signer.Sign(content)
	if err != nil {
		return nil, fmt.Errorf("failed to sign content: %w", err)
	}
	return signed, nil
}
