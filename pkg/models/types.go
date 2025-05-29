package models

// W3C Verifiable Credentials contexts
const (
	W3CCredentialsContext = "https://www.w3.org/2018/credentials/v1"
	ACKIDContext          = "https://agentcommercekit.com/contexts/ack-id/v1"
)

// Standard VC types
var (
	VerifiableCredentialType        = []string{"VerifiableCredential"}
	AgentAuthorizationCredentialType = []string{"VerifiableCredential", "AgentAuthorizationCredential"}
	AgentOwnershipCredentialType     = []string{"VerifiableCredential", "AgentOwnershipCredential"}
)

// Standard contexts
var (
	StandardContexts = []string{W3CCredentialsContext, ACKIDContext}
)

// Common action types for agent authorization
const (
	ActionTransfer         = "transfer"
	ActionBooking          = "booking"
	ActionQuote            = "quote"
	ActionPayment          = "payment"
	ActionMessage          = "message"
	ActionRead             = "read"
	ActionWrite            = "write"
	ActionExecute          = "execute"
)

// Common scope/resource types
const (
	ScopeETH           = "ETH"
	ScopeUSD           = "USD"
	ScopeEUR           = "EUR"
	ScopeFlights       = "flights"
	ScopeHotels        = "hotels"
	ScopeEmail         = "email"
	ScopeDatabase      = "database"
	ScopeAPI           = "api"
)

// DID method prefixes
const (
	DIDMethodAckid = "did:ackid:"
	DIDMethodWeb    = "did:web:"
	DIDMethodKey    = "did:key:"
)

type ClaimStatus string

const (
	StatusActive    ClaimStatus = "active"
	StatusRevoked   ClaimStatus = "revoked"
	StatusExpired   ClaimStatus = "expired"
	StatusSuspended ClaimStatus = "suspended"
)

type RequestPriority string

const (
	PriorityLow    RequestPriority = "low"
	PriorityNormal RequestPriority = "normal"
	PriorityHigh   RequestPriority = "high"
	PriorityUrgent RequestPriority = "urgent"
)

type RequestContext struct {
	SessionID       string                 `json:"session_id,omitempty"`        // For multi-step operations
	ParentRequestID string                 `json:"parent_request_id,omitempty"` // For related requests
	Environment     string                 `json:"environment,omitempty"`       // "production", "testing", "staging"
	IPAddress       string                 `json:"ip_address,omitempty"`        // For security monitoring
	UserAgent       string                 `json:"user_agent,omitempty"`        // For client identification
	CustomData      map[string]interface{} `json:"custom_data,omitempty"`       // For specific use cases
}

// EIP712Domain defines the domain for EIP-712 typed data signing
// This makes signatures more secure and human-readable in wallets
type EIP712Domain struct {
	Name              string `json:"name"`              // e.g., "EVMkya"
	Version           string `json:"version"`           // e.g., "1"
	ChainID           int64  `json:"chainId"`           // Ethereum chain ID
	VerifyingContract string `json:"verifyingContract,omitempty"` // Contract address if applicable
}