package server

import (
	"encoding/json"
	"net/http"
	"time"

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
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.NDC == "" || req.Quantity <= 0 || req.NPI == "" {
		writeError(w, http.StatusBadRequest, "ndc, npi and quantity are required")
		return
	}

	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	if req.Price < 0 {
		writeError(w, http.StatusBadRequest, "Price cannot be negative")
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

	response := APIResponse{
		Success: true,
		Message: "Claim created successfully",
		Data:    convertDBClaimToAPI(claim),
	}

	writeJSON(w, http.StatusCreated, response)
}

/*
// getClaim handles GET /api/v1/claims/{id}
func (server *Server) getClaim(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		writeError(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	id := pathParts[4]

	// Get claim from database
	claim, err := server.store.GetClaim(r.Context(), id)
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
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.ClaimID == "" {
		writeError(w, http.StatusBadRequest, "Claim ID is required")
		return
	}

	// Create reversal in database
	reversal, err := server.store.CreateReversal(r.Context(), req.ClaimID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create reversal")
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Reversal created successfully",
		Data:    convertDBReversalToAPI(reversal),
	}

	writeJSON(w, http.StatusCreated, response)
}
*/
