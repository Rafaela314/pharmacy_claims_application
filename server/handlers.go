package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	sqlc "github.com/pharmacy_claims_application/db/sqlc"
)

// healthCheck handles the health check endpoint
func (server *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Message: "Server is healthy",
		Data: map[string]interface{}{
			"timestamp": time.Now().UTC(),
			"status":    "ok",
		},
	}

	writeJSON(w, http.StatusOK, response)
}

// createClaim handles POST /api/v1/claims
func (server *Server) createClaim(w http.ResponseWriter, r *http.Request) {

	var req CreateClaimRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format in request body", map[string]interface{}{
			"expected_format": "JSON object with fields: ndc (string), npi (string), quantity (integer), price (number)",
			"example": map[string]interface{}{
				"ndc":      "123456789",
				"npi":      "9876543210",
				"quantity": 30,
				"price":    15.99,
			},
		})
		return
	}

	// Basic validation with specific error messages
	if req.NDC == "" {
		writeError(w, http.StatusBadRequest, "NDC (National Drug Code) is required", map[string]interface{}{
			"field":       "ndc",
			"type":        "string",
			"description": "National Drug Code identifier",
			"example":     "123456789",
		})
		return
	}

	if req.NPI == "" {
		writeError(w, http.StatusBadRequest, "NPI (National Provider Identifier) is required", map[string]interface{}{
			"field":       "npi",
			"type":        "string",
			"description": "National Provider Identifier",
			"example":     "9876543210",
		})
		return
	}

	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "Quantity must be greater than 0", map[string]interface{}{
			"field":     "quantity",
			"type":      "integer",
			"min_value": 1,
			"example":   30,
		})
		return
	}

	if req.Price < 0 {
		writeError(w, http.StatusBadRequest, "Price cannot be negative", map[string]interface{}{
			"field":     "price",
			"type":      "number",
			"min_value": 0,
			"example":   15.99,
		})
		return
	}

	// Create claim in database
	arg := sqlc.CreateClaimParams{
		NDC:      req.NDC,
		NPI:      req.NPI,
		Quantity: int64(req.Quantity),
		Price:    req.Price,
	}

	claim, err := server.store.CreateClaim(r.Context(), arg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create claim")
		return
	}

	// Log the claim submission event
	if err := server.logger.LogClaimSubmission(claim.ID, req.NDC, req.NPI, req.Quantity, req.Price); err != nil {
		log.Printf("Warning: failed to log claim submission: %v", err)
	}

	response := map[string]interface{}{
		"status":   "claim submitted",
		"claim_id": claim.ID.String(),
	}

	writeJSON(w, http.StatusCreated, response)
}

// getClaim handles GET /api/v1/claims/{id}
func (server *Server) getClaim(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		writeError(w, http.StatusBadRequest, "Invalid URL format. Expected: /api/v1/claims/{id}")
		return
	}

	id := pathParts[4]
	if id == "" {
		writeError(w, http.StatusBadRequest, "Claim ID cannot be empty")
		return
	}

	// Parse UUID
	claimID, err := uuid.Parse(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid claim ID format. Must be a valid UUID", map[string]interface{}{
			"field":   "claim_id",
			"type":    "string (UUID)",
			"format":  "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"example": "550e8400-e29b-41d4-a716-446655440000",
		})
		return
	}

	// Get claim from database
	claim, err := server.store.GetClaim(r.Context(), claimID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Claim not found")
		return
	}

	response := APIResponse{
		Success: true,
		Data:    convertDBClaimToAPI(claim),
	}

	writeJSON(w, http.StatusOK, response)
}

// createReversal handles POST /api/v1/reversals
func (server *Server) createReversal(w http.ResponseWriter, r *http.Request) {
	var req CreateReversalRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format in request body", map[string]interface{}{
			"expected_format": "JSON object with field: claim_id (string, UUID format)",
			"example": map[string]interface{}{
				"claim_id": "550e8400-e29b-41d4-a716-446655440000",
			},
		})
		return
	}

	// Basic validation
	if req.ClaimID == uuid.Nil {
		writeError(w, http.StatusBadRequest, "Claim ID is required and must be a valid UUID", map[string]interface{}{
			"field":   "claim_id",
			"type":    "string (UUID)",
			"format":  "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"example": "550e8400-e29b-41d4-a716-446655440000",
		})
		return
	}

	// Create reversal in database
	_, err := server.store.CreateReversal(r.Context(), req.ClaimID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create reversal")
		return
	}

	// Log the claim reversal event
	if err := server.logger.LogClaimReversal(req.ClaimID); err != nil {
		log.Printf("Warning: failed to log claim reversal: %v", err)
	}

	response := map[string]interface{}{
		"status":   "claim reversed",
		"claim_id": req.ClaimID.String(),
	}

	writeJSON(w, http.StatusCreated, response)
}
