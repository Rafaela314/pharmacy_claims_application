package db

import (
	"context"
	"testing"

	"github.com/pharmacy_claims_application/util"
	"github.com/stretchr/testify/require"
)

func createRandomPharmacy(t *testing.T) Pharmacy {
	arg := CreatePharmacyParams{
		NPI:   util.RandomNumericString(10),
		Chain: util.RandomString(10),
	}
	pharmacy, err := testQueries.CreatePharmacy(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, pharmacy)

	require.Equal(t, arg.NPI, pharmacy.NPI)
	require.Equal(t, arg.Chain, pharmacy.Chain)

	return pharmacy
}

func TestCreatePharmacy(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
		pharmacyArg := CreatePharmacyParams{
			NPI:   util.RandomNumericString(10),
			Chain: util.RandomString(10),
		}
		pharmacy, err := txQueries.CreatePharmacy(context.Background(), pharmacyArg)
		require.NoError(t, err)
		require.NotEmpty(t, pharmacy)

		require.Equal(t, pharmacyArg.NPI, pharmacy.NPI)
		require.Equal(t, pharmacyArg.Chain, pharmacy.Chain)
	})
}

func TestGetPharmacy(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
		// Create pharmacy within transaction
		pharmacyArg := CreatePharmacyParams{
			NPI:   util.RandomNumericString(10),
			Chain: util.RandomString(10),
		}
		pharmacy1, err := txQueries.CreatePharmacy(context.Background(), pharmacyArg)
		require.NoError(t, err)
		require.NotEmpty(t, pharmacy1)

		// Test retrieval within the same transaction
		pharmacy2, err := txQueries.GetPharmacy(context.Background(), pharmacy1.NPI)
		require.NoError(t, err)
		require.NotEmpty(t, pharmacy2)

		require.Equal(t, pharmacy1.NPI, pharmacy2.NPI)
		require.Equal(t, pharmacy1.Chain, pharmacy2.Chain)
	})
}
