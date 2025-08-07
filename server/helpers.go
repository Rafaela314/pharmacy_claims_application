package server

import (
	"encoding/json"
	"net/http"
	"time"

	sqlc "github.com/pharmacy_claims_application/db/sqlc"
)

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// writeError writes an error response
func writeError(w http.ResponseWriter, statusCode int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}

	writeJSON(w, statusCode, response)
}

// convertDBClaimToAPI converts a database claim to API format
func convertDBClaimToAPI(dbClaim sqlc.Claim) Claim {
	return Claim{
		ID:        dbClaim.ID.String(),
		NDC:       dbClaim.NDC,
		Quantity:  int(dbClaim.Quantity),
		NPI:       dbClaim.NPI,
		Price:     dbClaim.Price,
		Timestamp: dbClaim.Timestamp,
	}
}

// convertDBReversalToAPI converts a database reversal to API format
func convertDBReversalToAPI(dbReversal sqlc.Reversal) Reversal {
	return Reversal{
		ID:        dbReversal.ID.String(),
		ClaimID:   dbReversal.ClaimID.String(),
		Timestamp: dbReversal.Timestamp,
	}
}

// parseTime parses a time string in RFC3339 format
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// formatTime formats a time.Time to RFC3339 string
func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
