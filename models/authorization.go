package models

// AuthorizationRequest represents a request to verify an agent's authority
// Used in agent-to-agent verification flows
type AuthorizationRequest struct {
	AgentDID      string        `json:"agent_did"`        // Agent requesting authorization
	TargetAction  string        `json:"target_action"`    // Action being requested
	TargetScope   string        `json:"target_scope"`     // Resource/scope for the action
	Amount        string        `json:"amount,omitempty"` // Amount if financial
	RequesterDID  string        `json:"requester_did"`    // Who is asking for verification
	Timestamp     int64         `json:"timestamp"`
	Nonce         string        `json:"nonce"`
	RequestID     string        `json:"request_id"`       // Unique request identifier
	Priority      RequestPriority `json:"priority"` // Request priority level
	Context       RequestContext  `json:"context"`  // Request context
}

// AuthorizationResponse represents the response to an authorization request
type AuthorizationResponse struct {
	Authorized      bool                   `json:"authorized"`
	Reason          string                 `json:"reason,omitempty"`
	ValidUntil      int64                  `json:"valid_until,omitempty"`
	RemainingBudget string                 `json:"remaining_budget,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Timestamp       int64                  `json:"timestamp"`
	ResponderDID    string                 `json:"responder_did"`
	Signature       string                 `json:"signature,omitempty"`
}