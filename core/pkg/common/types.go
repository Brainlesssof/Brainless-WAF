package common

import (
	"net/http"
	"net/url"
)

// Transaction wraps an HTTP request with its security context and normalized data.
type Transaction struct {
	ID             string
	Request        *http.Request
	NormalizedURL  string
	NormalizedPath string
	Args           url.Values
	Headers        http.Header
	Body           []byte
	ContentType    string

	// Security metadata
	AnomalyScore int
	Variables    map[string]string // Transactional variables (TX:name)
	Action       string            // deny, allow, log, pass
	MatchedRules []int             // Rule IDs
}

func NewTransaction(r *http.Request) *Transaction {
	id := r.Header.Get("X-Request-ID")
	return &Transaction{
		ID:        id,
		Request:   r,
		Headers:   make(http.Header),
		Args:      make(url.Values),
		Variables: make(map[string]string),
		Action:    "allow", // Default to allow
	}
}
