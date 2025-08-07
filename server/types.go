package server

import (
	"time"

	"github.com/google/uuid"
)

// Claim represents a pharmacy claim
type Claim struct {
	ID        string    `json:"id"`
	NDC       string    `json:"ndc"`
	Quantity  int       `json:"quantity"`
	NPI       string    `json:"npi"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

// Reversal represents a pharmacy claim reversal
type Reversal struct {
	ID        string    `json:"id"`
	ClaimID   string    `json:"claim_id"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateClaimRequest represents the request body for creating a claim
type CreateClaimRequest struct {
	NDC      string  `json:"ndc" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	NPI      string  `json:"npi" validate:"required"`
	Price    float64 `json:"price" validate:"required,min=0"`
}

// CreateReversalRequest represents the request body for creating a reversal
type CreateReversalRequest struct {
	ClaimID uuid.UUID `json:"claim_id" validate:"required"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  string      `json:"status,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
