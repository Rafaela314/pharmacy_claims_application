package db

import (
	"context"
	"testing"
	"time"

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
	createRandomReversal(t)
}

func TestGetReversalByClaimID(t *testing.T) {
	reversal1 := createRandomReversal(t)
	reversal2, err := testQueries.GetReversalByClaimID(context.Background(), reversal1.ClaimID)
	require.NoError(t, err)
	require.NotEmpty(t, reversal2)

	require.Equal(t, reversal1.ClaimID, reversal2.ClaimID)
	require.WithinDuration(t, reversal1.Timestamp, reversal2.Timestamp, time.Second)
}
