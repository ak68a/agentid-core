package models

import (
	"fmt"
	"time"
)

// DelegationClaim represents a delegation from one agent to another
// This allows agents to sub-delegate authority within defined constraints
type DelegationClaim struct {
	// Core Delegation Info
	DelegatorDID string `json:"delegator_did"` // Agent granting authority
	DelegateDID  string `json:"delegate_did"`  // Agent receiving authority
	
	// What's being delegated
	Action      string                 `json:"action"`      // e.g., "transfer", "book_flight"
	Scope       string                 `json:"scope"`       // e.g., "ETH", "hotels"
	Constraints map[string]interface{} `json:"constraints"` // Time limits, scope restrictions, etc.
	
	// Temporal bounds
	IssuedAt  int64 `json:"issued_at"`
	ExpiresAt int64 `json:"expires_at"`
	Nonce     string `json:"nonce"`
	
	// Chain tracking
	ParentDelegation *string `json:"parent_delegation,omitempty"` // Reference to parent in chain
	MaxDepth         int     `json:"max_depth"`                   // How many levels can be sub-delegated
	CurrentDepth     int     `json:"current_depth"`               // Current position in chain
	
	// W3C VC compliance
	Type    []string `json:"type,omitempty"`
	Context []string `json:"@context,omitempty"`
	Issuer  string   `json:"issuer,omitempty"`
	Subject string   `json:"subject,omitempty"`
	Proof   *CredentialProof `json:"proof,omitempty"`
}

// DelegationChain represents a complete chain of delegations
type DelegationChain struct {
	Delegations []*DelegationClaim `json:"delegations"`
	Valid       bool               `json:"valid"`
	Reason      string             `json:"reason,omitempty"`
}

// TimeConstraint for temporal limitations
type TimeConstraint struct {
	ValidFrom  int64   `json:"valid_from"`            // Unix timestamp
	ValidUntil int64   `json:"valid_until"`           // Unix timestamp
	TimeZone   string  `json:"timezone,omitempty"`    // e.g., "UTC", "America/New_York"
	Days       []int   `json:"days,omitempty"`        // Days of week (0=Sunday)
	Hours      []int   `json:"hours,omitempty"`       // Hours of day (0-23)
}

// ScopeConstraint for limiting what resources can be accessed
type ScopeConstraint struct {
	AllowedResources []string               `json:"allowed_resources"`
	DeniedResources  []string               `json:"denied_resources,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// ChainError represents errors that can occur during chain operations
type ChainError struct {
	Code    string
	Message string
}

func (e *ChainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// IsExpired checks if a delegation has expired
func (dc *DelegationClaim) IsExpired() bool {
	if dc.ExpiresAt == 0 {
		return false // Never expires
	}
	return time.Now().Unix() > dc.ExpiresAt
}

// CanSubDelegate checks if this delegation allows further sub-delegation
func (dc *DelegationClaim) CanSubDelegate() bool {
	return dc.CurrentDepth < dc.MaxDepth
}

// GetTimeConstraint extracts time constraint from the delegation
func (dc *DelegationClaim) GetTimeConstraint() *TimeConstraint {
	if timeConstr, ok := dc.Constraints["time"].(map[string]interface{}); ok {
		return &TimeConstraint{
			ValidFrom:  getInt64(timeConstr, "valid_from"),
			ValidUntil: getInt64(timeConstr, "valid_until"),
			TimeZone:   getString(timeConstr, "timezone"),
			Days:       getIntSlice(timeConstr, "days"),
			Hours:      getIntSlice(timeConstr, "hours"),
		}
	}
	return nil
}

// ValidateChain validates an entire delegation chain
func (chain *DelegationChain) ValidateChain() bool {
	if len(chain.Delegations) == 0 {
		chain.Valid = false
		chain.Reason = "empty delegation chain"
		return false
	}

	// Check each delegation in the chain
	for i, delegation := range chain.Delegations {
		// Check expiration
		if delegation.IsExpired() {
			chain.Valid = false
			chain.Reason = fmt.Sprintf("delegation %d is expired", i)
			return false
		}

		// Check depth constraints
		if delegation.CurrentDepth > delegation.MaxDepth {
			chain.Valid = false
			chain.Reason = fmt.Sprintf("delegation %d exceeds max depth", i)
			return false
		}

		// Check chain continuity (delegator of next should be delegate of previous)
		if i > 0 {
			prevDelegation := chain.Delegations[i-1]
			if prevDelegation.DelegateDID != delegation.DelegatorDID {
				chain.Valid = false
				chain.Reason = fmt.Sprintf("broken chain at delegation %d", i)
				return false
			}
		}
	}

	chain.Valid = true
	return true
}

// Helper functions for extracting values from constraint maps
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getInt64(m map[string]interface{}, key string) int64 {
	if val, ok := m[key].(float64); ok {
		return int64(val)
	}
	if val, ok := m[key].(int64); ok {
		return val
	}
	return 0
}

func getIntSlice(m map[string]interface{}, key string) []int {
	if val, ok := m[key].([]interface{}); ok {
		result := make([]int, len(val))
		for i, v := range val {
			if num, ok := v.(float64); ok {
				result[i] = int(num)
			}
		}
		return result
	}
	return nil
}

// GetChain builds and returns the complete delegation chain for a claim
// by following parent references
func (dc *DelegationClaim) GetChain() (*DelegationChain, error) {
	chain := &DelegationChain{
		Delegations: []*DelegationClaim{dc},
	}

	// Follow parent references to build the chain
	current := dc
	for current.ParentDelegation != nil {
		// In a real implementation, you would fetch the parent delegation
		// from storage using the ParentDelegation reference
		// For now, we'll return an error
		return nil, &ChainError{
			Code:    "PARENT_NOT_FOUND",
			Message: "parent delegation reference not implemented",
		}
	}

	// Validate the chain
	if !chain.ValidateChain() {
		return nil, &ChainError{
			Code:    "INVALID_CHAIN",
			Message: chain.Reason,
		}
	}

	return chain, nil
}

// ValidateInChain validates this delegation in the context of its chain
func (dc *DelegationClaim) ValidateInChain() (*DelegationChain, error) {
	chain, err := dc.GetChain()
	if err != nil {
		return nil, err
	}

	// Additional chain-specific validations
	if err := chain.ValidateChainConstraints(); err != nil {
		return nil, err
	}

	return chain, nil
}

// GetRootDelegation returns the original delegator in the chain
func (chain *DelegationChain) GetRootDelegation() *DelegationClaim {
	if len(chain.Delegations) == 0 {
		return nil
	}
	return chain.Delegations[0]
}

// GetLeafDelegation returns the final delegate in the chain
func (chain *DelegationChain) GetLeafDelegation() *DelegationClaim {
	if len(chain.Delegations) == 0 {
		return nil
	}
	return chain.Delegations[len(chain.Delegations)-1]
}

// ValidateChainConstraints performs chain-wide constraint validation
func (chain *DelegationChain) ValidateChainConstraints() error {
	if len(chain.Delegations) == 0 {
		return &ChainError{
			Code:    "EMPTY_CHAIN",
			Message: "delegation chain is empty",
		}
	}

	// Track chain-wide constraints
	var earliestValidFrom int64
	var latestValidUntil int64
	var allowedResources map[string]bool
	var deniedResources map[string]bool

	// Initialize constraints from first delegation
	root := chain.GetRootDelegation()
	if timeConstr := root.GetTimeConstraint(); timeConstr != nil {
		earliestValidFrom = timeConstr.ValidFrom
		latestValidUntil = timeConstr.ValidUntil
	}
	if scopeConstr := root.GetScopeConstraint(); scopeConstr != nil {
		allowedResources = make(map[string]bool)
		deniedResources = make(map[string]bool)
		for _, resource := range scopeConstr.AllowedResources {
			allowedResources[resource] = true
		}
		for _, resource := range scopeConstr.DeniedResources {
			deniedResources[resource] = true
		}
	}

	// Validate each delegation in the chain
	for i, delegation := range chain.Delegations {
		// Skip root delegation as we used it for initialization
		if i == 0 {
			continue
		}

		// Validate time constraints
		if timeConstr := delegation.GetTimeConstraint(); timeConstr != nil {
			if timeConstr.ValidFrom < earliestValidFrom {
				return &ChainError{
					Code:    "TIME_CONSTRAINT_VIOLATION",
					Message: fmt.Sprintf("delegation %d starts before chain start time", i),
				}
			}
			if timeConstr.ValidUntil > latestValidUntil {
				return &ChainError{
					Code:    "TIME_CONSTRAINT_VIOLATION",
					Message: fmt.Sprintf("delegation %d ends after chain end time", i),
				}
			}
		}

		// Validate scope constraints
		if scopeConstr := delegation.GetScopeConstraint(); scopeConstr != nil {
			for _, resource := range scopeConstr.AllowedResources {
				if deniedResources[resource] {
					return &ChainError{
						Code:    "SCOPE_CONSTRAINT_VIOLATION",
						Message: fmt.Sprintf("delegation %d allows denied resource: %s", i, resource),
					}
				}
			}
			for _, resource := range scopeConstr.DeniedResources {
				if allowedResources[resource] {
					return &ChainError{
						Code:    "SCOPE_CONSTRAINT_VIOLATION",
						Message: fmt.Sprintf("delegation %d denies allowed resource: %s", i, resource),
					}
				}
			}
		}
	}

	return nil
}

// RevokeChain generates revocation claims for the entire chain
func (chain *DelegationChain) RevokeChain(reason string) []*RevocationClaim {
	var revocations []*RevocationClaim
	now := time.Now().Unix()

	for _, delegation := range chain.Delegations {
		// Create a revocation claim for each delegation
		revocation := &RevocationClaim{
			RevokedCredentialID: delegation.Nonce, // Using nonce as credential ID
			RevokedAgentDID:     delegation.DelegateDID,
			RevokerDID:          chain.GetRootDelegation().DelegatorDID, // Root delegator revokes
			Reason:              reason,
			RevokedAt:           now,
			EffectiveAt:         now, // Immediate effect
			Nonce:              fmt.Sprintf("revoke_%d_%s", now, delegation.Nonce),
			Type:               RevocationCredentialType,
			Context:            StandardContexts,
			Issuer:             chain.GetRootDelegation().DelegatorDID,
			Subject:            delegation.DelegateDID,
		}
		revocations = append(revocations, revocation)
	}

	return revocations
}

// GetScopeConstraint extracts scope constraint from the delegation
func (dc *DelegationClaim) GetScopeConstraint() *ScopeConstraint {
	if scope, ok := dc.Constraints["scope"].(map[string]interface{}); ok {
		return &ScopeConstraint{
			AllowedResources: getStringSlice(scope, "allowed_resources"),
			DeniedResources:  getStringSlice(scope, "denied_resources"),
			Metadata:         getMap(scope, "metadata"),
		}
	}
	return nil
}

// Helper function to get string slice from interface map
func getStringSlice(m map[string]interface{}, key string) []string {
	if val, ok := m[key].([]interface{}); ok {
		result := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				result[i] = str
			}
		}
		return result
	}
	return nil
}

// Helper function to get map from interface map
func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if val, ok := m[key].(map[string]interface{}); ok {
		return val
	}
	return nil
}