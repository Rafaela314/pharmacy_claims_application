package db

import (
	"context"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pharmacy_claims_application/util"
	"github.com/stretchr/testify/require"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../../pharmacy_claims_application")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testDB = conn
	testQueries = New(conn)

	exitCode := m.Run()

	if testDB != nil {
		testDB.Close()
	}

	os.Exit(exitCode)
}

// runTestWithTransaction runs a test function within a transaction that gets rolled back
func runTestWithTransaction(t *testing.T, testFunc func(*testing.T, *Queries)) {
	ctx := context.Background()

	tx, err := testDB.Begin(ctx)
	require.NoError(t, err)

	txQueries := New(tx)

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			t.Logf("Warning: failed to rollback transaction: %v", err)
		}
	}()

	// Run the test with transaction queries
	testFunc(t, txQueries)
}

// createTestData creates sample data for testing
func createTestData(t *testing.T) (Pharmacy, Claim, Reversal) {

	pharmacy := createRandomPharmacy(t)

	claim := createRandomClaimWithPharmacy(t, pharmacy)

	reversal := createRandomReversalWithClaim(t, claim)

	return pharmacy, claim, reversal
}

// createRandomClaimWithPharmacy creates a claim using the provided pharmacy
func createRandomClaimWithPharmacy(t *testing.T, pharmacy Pharmacy) Claim {
	arg := CreateClaimParams{
		NDC:      util.RandomString(11),
		Price:    util.RandomMoney(),
		Quantity: util.RandomInt(1, 1000),
		NPI:      pharmacy.NPI,
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

// createRandomReversalWithClaim creates a reversal using the provided claim
func createRandomReversalWithClaim(t *testing.T, claim Claim) Reversal {
	reversal, err := testQueries.CreateReversal(context.Background(), claim.ID)
	require.NoError(t, err)
	require.NotEmpty(t, reversal)

	require.Equal(t, claim.ID, reversal.ClaimID)
	require.NotZero(t, reversal.ID)
	require.NotZero(t, reversal.Timestamp)

	return reversal
}
