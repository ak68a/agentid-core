package models

import "time"

// RevocationClaim represents a revocation of a previously issued credential
type RevocationClaim struct {
	// What's being revoked
	RevokedCredentialID string `json:"revoked_credential_id"` // Hash or ID of the revoked credential
	RevokedAgentDID     string `json:"revoked_agent_did"`     // Agent whose credential is revoked
	
	// Who's doing the revocation
	RevokerDID string `json:"revoker_did"` // Must be the original issuer or have revocation authority
	
	// Why and when
	Reason       string `json:"reason"`                    // "compromised", "expired", "policy_change", etc.
	RevokedAt    int64  `json:"revoked_at"`               // When revocation occurred
	EffectiveAt  int64  `json:"effective_at,omitempty"`   // When revocation takes effect (can be future)
	Nonce        string `json:"nonce"`
	
	// Optional details
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	
	// W3C VC compliance
	Type    []string `json:"type,omitempty"`
	Context []string `json:"@context,omitempty"`
	Issuer  string   `json:"issuer,omitempty"`
	Subject string   `json:"subject,omitempty"`
	Proof   *CredentialProof `json:"proof,omitempty"`
}

// RevocationList represents a collection of revoked credentials
type RevocationList struct {
	ListID      string             `json:"list_id"`      // Unique identifier for this revocation list
	IssuerDID   string             `json:"issuer_did"`   // Who maintains this list
	LastUpdated int64              `json:"last_updated"` // When list was last updated
	Revocations []*RevocationClaim `json:"revocations"`  // All revocations in this list
	
	// W3C compliance
	Type    []string `json:"type,omitempty"`
	Context []string `json:"@context,omitempty"`
	Proof   *CredentialProof `json:"proof,omitempty"`
}

// RevocationQuery represents a query to check if a credential is revoked
type RevocationQuery struct {
	CredentialID string `json:"credential_id"` // ID of credential to check
	AgentDID     string `json:"agent_did"`     // Agent DID to check
	AsOfTime     int64  `json:"as_of_time"`    // Check revocation status as of this time
	QueryID      string `json:"query_id"`      // Unique query identifier
	RequesterDID string `json:"requester_did"` // Who is asking
	Timestamp    int64  `json:"timestamp"`
	Nonce        string `json:"nonce"`
}

// RevocationStatus represents the response to a revocation query
type RevocationStatus struct {
	CredentialID string `json:"credential_id"`
	AgentDID     string `json:"agent_did"`
	IsRevoked    bool   `json:"is_revoked"`
	RevokedAt    int64  `json:"revoked_at,omitempty"`    // When it was revoked
	Reason       string `json:"reason,omitempty"`        // Why it was revoked
	QueryID      string `json:"query_id"`                // Matches the query
	ResponderDID string `json:"responder_did"`           // Who provided this response
	Timestamp    int64  `json:"timestamp"`
	Signature    string `json:"signature,omitempty"`
}

// Revocation reasons (constants)
const (
	RevocationReasonCompromised   = "compromised"
	RevocationReasonExpired       = "expired"
	RevocationReasonPolicyChange  = "policy_change"
	RevocationReasonSuspension    = "suspension"
	RevocationReasonTermination   = "termination"
	RevocationReasonKeyRotation   = "key_rotation"
	RevocationReasonUnauthorized  = "unauthorized"
)

// VC types for revocation
var (
	RevocationCredentialType = []string{"VerifiableCredential", "RevocationCredential"}
	RevocationListType       = []string{"VerifiableCredential", "RevocationList2020"}
)

// IsEffective checks if a revocation is currently in effect
func (rc *RevocationClaim) IsEffective() bool {
	now := time.Now().Unix()
	if rc.EffectiveAt == 0 {
		// No specific effective time, use revoked time
		return now >= rc.RevokedAt
	}
	return now >= rc.EffectiveAt
}

// IsRevoked checks if a specific credential ID is in the revocation list
func (rl *RevocationList) IsRevoked(credentialID string) *RevocationClaim {
	for _, revocation := range rl.Revocations {
		if revocation.RevokedCredentialID == credentialID && revocation.IsEffective() {
			return revocation
		}
	}
	return nil
}

// IsAgentRevoked checks if any credentials for an agent DID are revoked
func (rl *RevocationList) IsAgentRevoked(agentDID string) []*RevocationClaim {
	var revocations []*RevocationClaim
	for _, revocation := range rl.Revocations {
		if revocation.RevokedAgentDID == agentDID && revocation.IsEffective() {
			revocations = append(revocations, revocation)
		}
	}
	return revocations
}

// AddRevocation adds a new revocation to the list
func (rl *RevocationList) AddRevocation(revocation *RevocationClaim) {
	rl.Revocations = append(rl.Revocations, revocation)
	rl.LastUpdated = time.Now().Unix()
}

// GetRevocationsSince returns all revocations added since a given timestamp
func (rl *RevocationList) GetRevocationsSince(since int64) []*RevocationClaim {
	var recent []*RevocationClaim
	for _, revocation := range rl.Revocations {
		if revocation.RevokedAt >= since {
			recent = append(recent, revocation)
		}
	}
	return recent
}