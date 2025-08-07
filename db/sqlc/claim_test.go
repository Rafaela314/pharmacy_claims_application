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

func TestCreateClaim(t *testing.T) {
	createRandomClaim(t)
}

func TestGetClaim(t *testing.T) {
	claim1 := createRandomClaim(t)
	claim2, err := testQueries.GetClaim(context.Background(), claim1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, claim2)

	require.Equal(t, claim1.ID, claim2.ID)
	require.Equal(t, claim1.NDC, claim2.NDC)
	require.Equal(t, claim1.NPI, claim2.NPI)
	require.Equal(t, claim1.Price, claim2.Price)
	require.Equal(t, claim1.Quantity, claim2.Quantity)

	require.WithinDuration(t, claim1.Timestamp, claim2.Timestamp, time.Second)
}
