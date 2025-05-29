package models

import "time"

// NewAgentClaim creates a new agent claim with proper W3C VC structure
func NewAgentClaim(agentDID, ownerDID, action, scope string, expiresAt int64, nonce string) *AgentClaim {
	now := time.Now().Unix()
	
	return &AgentClaim{
		AgentDID:  agentDID,
		OwnerDID:  ownerDID,
		Action:    action,
		Scope:     scope,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
		Nonce:     nonce,
		Type:      AgentAuthorizationCredentialType,
		Context:   StandardContexts,
		Issuer:    ownerDID,
		Subject:   agentDID,
	}
}

// NewOwnershipClaim creates a new ownership claim
func NewOwnershipClaim(agentDID, ownerDID string, nonce string) *OwnershipClaim {
	now := time.Now().Unix()
	
	return &OwnershipClaim{
		AgentDID:  agentDID,
		OwnerDID:  ownerDID,
		IssuedAt:  now,
		ExpiresAt: 0, // Never expires by default
		Nonce:     nonce,
		Type:      AgentOwnershipCredentialType,
		Context:   StandardContexts,
		Issuer:    ownerDID,
		Subject:   agentDID,
	}
}

// NewTransferClaim creates an authorization claim for token transfers
func NewTransferClaim(agentDID, ownerDID, token, maxAmount string, expiresAt int64, nonce string) *AgentClaim {
	claim := NewAgentClaim(agentDID, ownerDID, ActionTransfer, token, expiresAt, nonce)
	claim.MaxAmount = maxAmount
	return claim
}

// NewQuotingClaim creates an authorization claim for providing quotes
func NewQuotingClaim(agentDID, ownerDID, currencyPair string, expiresAt int64, nonce string) *AgentClaim {
	return NewAgentClaim(agentDID, ownerDID, ActionQuote, currencyPair, expiresAt, nonce)
}

// NewBookingClaim creates an authorization claim for making bookings
func NewBookingClaim(agentDID, ownerDID, bookingType string, expiresAt int64, nonce string) *AgentClaim {
	return NewAgentClaim(agentDID, ownerDID, ActionBooking, bookingType, expiresAt, nonce)
}

// NewAuthorizationRequest creates a new authorization request
func NewAuthorizationRequest(agentDID, targetAction, targetScope, requesterDID, nonce string) *AuthorizationRequest {
	return &AuthorizationRequest{
		AgentDID:     agentDID,
		TargetAction: targetAction,
		TargetScope:  targetScope,
		RequesterDID: requesterDID,
		Timestamp:    time.Now().Unix(),
		Nonce:        nonce,
	}
}

// NewAuthorizationResponse creates a new authorization response
func NewAuthorizationResponse(authorized bool, reason, responderDID string) *AuthorizationResponse {
	return &AuthorizationResponse{
		Authorized:   authorized,
		Reason:       reason,
		Timestamp:    time.Now().Unix(),
		ResponderDID: responderDID,
	}
}

// NewCredentialProof creates a new credential proof with EIP-712 domain
func NewCredentialProof(suite ProofSuite, purpose ProofPurpose) *CredentialProof {
	return &CredentialProof{
		Type:     string(suite),
		Created:  time.Now().Format(time.RFC3339),
		ProofPurpose: string(purpose),
		Domain: EIP712Domain{
			Name:    "AgentID",
			Version: "1",
			ChainID: 1, // Mainnet
		},
	}
}