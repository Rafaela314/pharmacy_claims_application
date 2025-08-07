package db

import (
	"context"
	"testing"
	"time"

	"github.com/pharmacy_claims_application/util"
	"github.com/stretchr/testify/require"
)

func createRandomReversal(t *testing.T) Reversal {
	// First create a claim
	claim := createRandomClaim(t)

	reversal, err := testQueries.CreateReversal(context.Background(), claim.ID)
	require.NoError(t, err)
	require.NotEmpty(t, reversal)

	require.Equal(t, claim.ID, reversal.ClaimID)

	require.NotZero(t, reversal.ID)
	require.NotZero(t, reversal.Timestamp)

	return reversal
}

func TestCreateReversal(t *testing.T) {
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

		// Create reversal within transaction
		reversal, err := txQueries.CreateReversal(context.Background(), claim.ID)
		require.NoError(t, err)
		require.NotEmpty(t, reversal)

		require.Equal(t, claim.ID, reversal.ClaimID)
		require.NotZero(t, reversal.ID)
		require.NotZero(t, reversal.Timestamp)
	})
}

func TestGetReversalByClaimID(t *testing.T) {
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

		// Create reversal within transaction
		reversal1, err := txQueries.CreateReversal(context.Background(), claim.ID)
		require.NoError(t, err)
		require.NotEmpty(t, reversal1)

		// Test retrieval within the same transaction
		reversal2, err := txQueries.GetReversalByClaimID(context.Background(), reversal1.ClaimID)
		require.NoError(t, err)
		require.NotEmpty(t, reversal2)

		require.Equal(t, reversal1.ClaimID, reversal2.ClaimID)
		require.WithinDuration(t, reversal1.Timestamp, reversal2.Timestamp, time.Second)
	})
}
