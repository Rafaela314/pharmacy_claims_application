package db

import (
	"context"
	"testing"

	"github.com/pharmacy_claims_application/util"
	"github.com/stretchr/testify/require"
)

// TestSimpleTransaction verifies that transaction-based testing works
func TestSimpleTransaction(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
		// Create test data within transaction
		pharmacyArg := CreatePharmacyParams{
			NPI:   util.RandomNumericString(10),
			Chain: util.RandomString(10),
		}
		pharmacy, err := txQueries.CreatePharmacy(context.Background(), pharmacyArg)
		require.NoError(t, err)
		require.NotEmpty(t, pharmacy.NPI)

		// Verify we can retrieve it within the same transaction
		retrievedPharmacy, err := txQueries.GetPharmacy(context.Background(), pharmacy.NPI)
		require.NoError(t, err)
		require.Equal(t, pharmacy.NPI, retrievedPharmacy.NPI)
	})
	// Transaction is rolled back, so no data should persist
}
