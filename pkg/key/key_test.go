package key

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestGenerate(t *testing.T) {
	agentKey, err := GenerateAgentKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	if agentKey.PrivateKey == nil {
		t.Error("Private key is nil")
	}

	if agentKey.PublicKey == nil {
		t.Error("Public key is nil")
	}

	if !strings.HasPrefix(agentKey.DID, "did:ackid:0x") {
		t.Errorf("DID format incorrect: %s", agentKey.DID)
	}

	if agentKey.Address.Hex() == "0x0000000000000000000000000000000000000000" {
		t.Error("Address is zero address")
	}
}

func TestImportFromHex(t *testing.T) {
	original, err := GenerateAgentKey()
	if err != nil {
		t.Fatalf("Failed to generate original key: %v", err)
	}

	hexKey := original.GetPrivateKeyHex()

	imported, err := ImportFromHex(hexKey)
	if err != nil {
		t.Fatalf("Failed to import key: %v", err)
	}

	if original.Address != imported.Address {
		t.Errorf("Addresses don't match: %s != %s", original.Address.Hex(), imported.Address.Hex())
	}

	if original.DID != imported.DID {
		t.Errorf("DIDs don't match: %s != %s", original.DID, imported.DID)
	}
}

func TestSignAndVerify(t *testing.T) {
	agentKey, err := GenerateAgentKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	message := []byte("Hello, Agent!")
	messageHash := crypto.Keccak256(message)

	signature, err := agentKey.Sign(messageHash)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	if len(signature) != 65 {
		t.Errorf("Signature length should be 65, got %d", len(signature))
	}

	// Verify with crypto.VerifySignature (using public key bytes)
	publicKeyBytes := crypto.FromECDSAPub(agentKey.PublicKey)
	isValid := crypto.VerifySignature(publicKeyBytes, messageHash, signature[:64])
	
	if !isValid {
		t.Error("Signature verification failed")
	}
}

func TestDIDFormat(t *testing.T) {
	agentKey, err := GenerateAgentKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	expectedPrefix := "did:ackid:"
	if !strings.HasPrefix(agentKey.DID, expectedPrefix) {
		t.Errorf("DID should start with %s, got %s", expectedPrefix, agentKey.DID)
	}

	// Should contain the address
	if !strings.Contains(agentKey.DID, agentKey.Address.Hex()) {
		t.Errorf("DID should contain address %s, got %s", agentKey.Address.Hex(), agentKey.DID)
	}
}