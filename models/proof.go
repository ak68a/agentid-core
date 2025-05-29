package models

// CredentialProof represents the cryptographic proof for a Verifiable Credential
// Following W3C VC Data Model specification
type CredentialProof struct {
	Type               string     `json:"type"`                          // e.g., "EcdsaSecp256k1Signature2019"
	Created            string     `json:"created"`                       // ISO 8601 timestamp
	VerificationMethod string     `json:"verificationMethod"`            // DID#key-id reference
	ProofPurpose       string     `json:"proofPurpose"`                  // e.g., "assertionMethod"
	ProofValue         string     `json:"proofValue"`                    // Base64 encoded signature
	Challenge          string     `json:"challenge,omitempty"`           // Random nonce for proof of possession
	Domain             EIP712Domain `json:"domain,omitempty"`            // EIP-712 domain for typed data signing
}

// ProofSuite defines the supported cryptographic proof types
type ProofSuite string

const (
	EcdsaSecp256k1Signature2019 ProofSuite = "EcdsaSecp256k1Signature2019"
	Ed25519Signature2018        ProofSuite = "Ed25519Signature2018"
	RsaSignature2018            ProofSuite = "RsaSignature2018"
)

// ProofPurpose defines the purpose of the cryptographic proof
type ProofPurpose string

const (
	AssertionMethod      ProofPurpose = "assertionMethod"
	Authentication       ProofPurpose = "authentication"
	KeyAgreement         ProofPurpose = "keyAgreement"
	CapabilityInvocation ProofPurpose = "capabilityInvocation"
	CapabilityDelegation ProofPurpose = "capabilityDelegation"
)