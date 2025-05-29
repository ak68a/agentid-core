package models

import "time"

// AgentClaim represents a verifiable claim about an agent's identity or authorization
// This aligns with W3C Verifiable Credentials standard used in ACK-ID
type AgentClaim struct {
	// Core Identity Claims
	AgentDID string `json:"agent_did"`            // The agent making the claim
	OwnerDID string `json:"owner_did,omitempty"`  // The owner who controls this agent

	// Lifecycle Management
	Status ClaimStatus `json:"status"`            // Current status of the claim
	
	// Action & Scope (what the agent is claiming authority to do)
	Action string `json:"action"`                 // e.g., "transfer", "book_flight", "provide_fx_quotes"
	Scope  string `json:"scope"`                  // e.g., "ETH", "USD", "hotel_bookings", resource being acted upon
	
	// Temporal & Security
	IssuedAt  int64  `json:"issued_at"`           // When claim was issued (Unix timestamp)
	ExpiresAt int64  `json:"expires_at"`          // When claim expires (Unix timestamp)
	Nonce     string `json:"nonce"`               // Prevent replay attacks
	
	// Optional Constraints
	MaxAmount string                 `json:"max_amount,omitempty"` // Maximum amount for financial actions
	Metadata  map[string]interface{} `json:"metadata,omitempty"`   // Additional constraints/context
	
	// W3C VC Standard Fields
	Type     []string `json:"type,omitempty"`     // e.g., ["VerifiableCredential", "AgentAuthorizationCredential"]
	Context  []string `json:"@context,omitempty"` // JSON-LD context for VC standard
	Issuer   string   `json:"issuer,omitempty"`   // DID of the issuer (usually the Owner)
	Subject  string   `json:"subject,omitempty"`  // DID of the subject (usually the Agent)
	
	// Cryptographic Proof
	Proof *CredentialProof `json:"proof,omitempty"`
}

// OwnershipClaim represents a claim that links an Agent to its Owner
// This is the foundational relationship in ACK-ID
type OwnershipClaim struct {
	AgentDID  string `json:"agent_did"`   // The agent being claimed
	OwnerDID  string `json:"owner_did"`   // The owner making the claim
	IssuedAt  int64  `json:"issued_at"`   // When ownership was established
	ExpiresAt int64  `json:"expires_at"`  // When ownership expires (0 = never)
	Nonce     string `json:"nonce"`
	
	// W3C VC compliance
	Type    []string `json:"type,omitempty"`
	Context []string `json:"@context,omitempty"`
	Issuer  string   `json:"issuer,omitempty"`
	Subject string   `json:"subject,omitempty"`
	Proof   *CredentialProof `json:"proof,omitempty"`
}

// IsExpired checks if a claim has expired
func (ac *AgentClaim) IsExpired() bool {
	if ac.ExpiresAt == 0 {
		return false // Never expires
	}
	return time.Now().Unix() > ac.ExpiresAt
}

// IsExpired checks if an ownership claim has expired
func (oc *OwnershipClaim) IsExpired() bool {
	if oc.ExpiresAt == 0 {
		return false // Never expires
	}
	return time.Now().Unix() > oc.ExpiresAt
}

// ToCredential converts an AgentClaim to a standard W3C VC format
func (ac *AgentClaim) ToCredential() map[string]interface{} {
	return map[string]interface{}{
		"@context": ac.Context,
		"type":     ac.Type,
		"issuer":   ac.Issuer,
		"issuanceDate":   time.Unix(ac.IssuedAt, 0).Format(time.RFC3339),
		"expirationDate": time.Unix(ac.ExpiresAt, 0).Format(time.RFC3339),
		"credentialSubject": map[string]interface{}{
			"id":         ac.Subject,
			"agentDID":   ac.AgentDID,
			"ownerDID":   ac.OwnerDID,
			"action":     ac.Action,
			"scope":      ac.Scope,
			"maxAmount":  ac.MaxAmount,
			"metadata":   ac.Metadata,
		},
		"proof": ac.Proof,
	}
}