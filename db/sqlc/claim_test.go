package db

import (
	"context"
	"testing"
	"time"

	"github.com/pharmacy_claims_application/util"
	"github.com/stretchr/testify/require"
)

func createRandomClaim(t *testing.T) Claim {
	// First create a pharmacy
	pharmacy := createRandomPharmacy(t)

	arg := CreateClaimParams{
		NDC:      util.RandomString(11),
		Price:    util.RandomMoney(),
		Quantity: util.RandomInt(1, 1000),
		NPI:      pharmacy.NPI, // Use the pharmacy's NPI
	}
	claim, err := testQueries.CreateClaim(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, claim)

	require.Equal(t, arg.NDC, claim.NDC)
	require.Equal(t, arg.Price, claim.Price)
	require.Equal(t, arg.NPI, claim.NPI)
	require.Equal(t, arg.Quantity, claim.Quantity)

	require.NotZero(t, claim.ID)
	require.NotZero(t, claim.Timestamp)

	return claim
}

func TestGetClaim(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
		// Create pharmacy within transaction
		pharmacyArg := CreatePharmacyParams{
			NPI:   util.RandomNumericString(10),
			Chain: util.RandomString(10),
		}
		pharmacy, err := txQueries.CreatePharmacy(context.Background(), pharmacyArg)
		require.NoError(t, err)
		require.NotEmpty(t, pharmacy)

		// Create claim within transaction
		claimArg := CreateClaimParams{
			NDC:      util.RandomString(11),
			Price:    util.RandomMoney(),
			Quantity: util.RandomInt(1, 1000),
			NPI:      pharmacy.NPI,
		}
		claim1, err := txQueries.CreateClaim(context.Background(), claimArg)
		require.NoError(t, err)
		require.NotEmpty(t, claim1)

		// Test retrieval within the same transaction
		claim2, err := txQueries.GetClaim(context.Background(), claim1.ID)
		require.NoError(t, err)
		require.NotEmpty(t, claim2)

		require.Equal(t, claim1.ID, claim2.ID)
		require.Equal(t, claim1.NDC, claim2.NDC)
		require.Equal(t, claim1.NPI, claim2.NPI)
		require.Equal(t, claim1.Price, claim2.Price)
		require.Equal(t, claim1.Quantity, claim2.Quantity)

		require.WithinDuration(t, claim1.Timestamp, claim2.Timestamp, time.Second)
	})
}

// TestCreateClaimWithTransaction demonstrates transaction-based testing
func TestCreateClaimWithTransaction(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
		// Create pharmacy within transaction
		pharmacyArg := CreatePharmacyParams{
			NPI:   util.RandomNumericString(10),
			Chain: util.RandomString(10),
		}
		pharmacy, err := txQueries.CreatePharmacy(context.Background(), pharmacyArg)
		require.NoError(t, err)
		require.NotEmpty(t, pharmacy)

		// Create claim within transaction
		claimArg := CreateClaimParams{
			NDC:      util.RandomString(11),
			Price:    util.RandomMoney(),
			Quantity: util.RandomInt(1, 1000),
			NPI:      pharmacy.NPI,
		}
		claim, err := txQueries.CreateClaim(context.Background(), claimArg)
		require.NoError(t, err)
		require.NotEmpty(t, claim)

		// Verify the claim was created correctly
		require.Equal(t, claimArg.NDC, claim.NDC)
		require.Equal(t, claimArg.Price, claim.Price)
		require.Equal(t, claimArg.NPI, claim.NPI)
		require.Equal(t, claimArg.Quantity, claim.Quantity)
		require.NotZero(t, claim.ID)
		require.NotZero(t, claim.Timestamp)

		// Test that we can retrieve the claim within the same transaction
		retrievedClaim, err := txQueries.GetClaim(context.Background(), claim.ID)
		require.NoError(t, err)
		require.Equal(t, claim.ID, retrievedClaim.ID)
		require.Equal(t, claim.NDC, retrievedClaim.NDC)
	})
	// Transaction is automatically rolled back here, so no data persists
}
