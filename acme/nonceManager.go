package acme

import (
	"encoding/json"
	"fmt"
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Manager Manages nonces.
type Manager struct {
	nonceURL string
	nonces   []string
	sync.Mutex
}

// NewManager Creates a new Manager.
func NewManager(nonceURL string) *Manager {
	return &Manager{
		nonceURL: nonceURL,
	}
}

// Pop Pops a nonce.
func (n *Manager) Pop() (string, bool) {
	n.Lock()
	defer n.Unlock()

	if len(n.nonces) == 0 {
		return "", false
	}

	nonce := n.nonces[len(n.nonces)-1]
	n.nonces = n.nonces[:len(n.nonces)-1]
	return nonce, true
}

// Push Pushes a nonce.
func (n *Manager) Push(nonce string) {
	n.Lock()
	defer n.Unlock()
	n.nonces = append(n.nonces, nonce)
}

// Nonce implement jose.NonceSource.
func (n *Manager) Nonce() (string, error) {
	if nonce, ok := n.Pop(); ok {
		return nonce, nil
	}
	return "", nil
}

// GetFromResponse Extracts a nonce from a MQTT message and pushes the nonce on the stack
func (n *Manager) GetAndPushFromResponse(msg MQTT.Message) {
	// parse replay-nonce out of message
	var nonce NewNonceResp
	err := json.Unmarshal(msg.Payload(), &nonce)
	if err != nil {
		fmt.Printf("Error parsing new-nonce response\n")
	}
	//push nonce on NonceManager
	n.Push(nonce.ReplayNonce)
}
