package key

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// AgentKey represents an agent's cryptographic identity
type AgentKey struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
	DID        string // DID format: did:ackid:0x{address}
}

// Generate creates a new random keypair for an agent
func GenerateAgentKey() (*AgentKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return fromPrivateKey(privateKey), nil
}

// ImportFromHex imports a private key from hexadecimal string
func ImportFromHex(hexKey string) (*AgentKey, error) {
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to import private key: %w", err)
	}

	return fromPrivateKey(privateKey), nil
}

// ImportFromBytes imports a private key from byte slice
func ImportFromBytes(keyBytes []byte) (*AgentKey, error) {
	privateKey, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to import private key: %w", err)
	}

	return fromPrivateKey(privateKey), nil
}


// Creates an AgentKey from an ECDSA private key
func fromPrivateKey(privateKey *ecdsa.PrivateKey) *AgentKey {
	publicKey := &privateKey.PublicKey
	address := crypto.PubkeyToAddress(*publicKey)
	did := fmt.Sprintf("did:ackid:%s", address.Hex())

	return &AgentKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		DID:        did,
	}
}

// Sign signs a message hash with the agent's private key
func (ak *AgentKey) Sign(messageHash []byte) ([]byte, error) {
	signature, err := crypto.Sign(messageHash, ak.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}
	return signature, nil
}

// GetPrivateKeyHex returns the private key as a hex string
func (ak *AgentKey) GetPrivateKeyHex() string {
	return fmt.Sprintf("%x", crypto.FromECDSA(ak.PrivateKey))
}

// GetPublicKeyHex returns the public key as a hex string  
func (ak *AgentKey) GetPublicKeyHex() string {
	return fmt.Sprintf("%x", crypto.FromECDSAPub(ak.PublicKey))
}

// VerifySignature verifies a signature against a message hash
func VerifySignature(publicKeyHex string, messageHash []byte, signature []byte) (bool, error) {
	// Remove recovery ID if present (last byte)
	if len(signature) == 65 {
		signature = signature[:64]
	}

	publicKey, err := crypto.HexToECDSA(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %w", err)
	}

	return crypto.VerifySignature(
		crypto.FromECDSAPub(&publicKey.PublicKey),
		messageHash,
		signature,
	), nil
}