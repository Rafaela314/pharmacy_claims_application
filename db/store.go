package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/pharmacy_claims_application/db/sqlc"
)

// Store defines all functions to execute database queries and transactions
type Store interface {
	CreateClaim(ctx context.Context, arg sqlc.CreateClaimParams) (sqlc.Claim, error)
	GetClaim(ctx context.Context, id uuid.UUID) (sqlc.Claim, error)
	CreateReversal(ctx context.Context, claimID uuid.UUID) (sqlc.Reversal, error)
	CreatePharmacy(ctx context.Context, arg sqlc.CreatePharmacyParams) (sqlc.Pharmacy, error)
	GetPharmacy(ctx context.Context, npi string) (sqlc.Pharmacy, error)
	CountPharmacies(ctx context.Context) (int64, error)
	CreateClaimTx(ctx context.Context, arg sqlc.CreateClaimParams) (sqlc.Claim, error)
	CreateReversalTx(ctx context.Context, claimID uuid.UUID) (sqlc.Reversal, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*sqlc.Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  sqlc.New(connPool),
	}
}

// CreateClaimTx creates a new claim within a database transaction
func (store *SQLStore) CreateClaimTx(ctx context.Context, arg sqlc.CreateClaimParams) (sqlc.Claim, error) {
	var result sqlc.Claim

	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		var err error

		result, err = q.CreateClaim(ctx, arg)
		return err
	})

	return result, err
}

// CreateReversalTx creates a new reversal within a database transaction
func (store *SQLStore) CreateReversalTx(ctx context.Context, claimID uuid.UUID) (sqlc.Reversal, error) {
	var result sqlc.Reversal

	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		var err error

		result, err = q.CreateReversal(ctx, claimID)
		return err
	})

	return result, err
}

// CreateClaim creates a new claim
func (store *SQLStore) CreateClaim(ctx context.Context, arg sqlc.CreateClaimParams) (sqlc.Claim, error) {
	return store.Queries.CreateClaim(ctx, arg)
}

// GetClaim gets a claim by ID
func (store *SQLStore) GetClaim(ctx context.Context, id uuid.UUID) (sqlc.Claim, error) {
	return store.Queries.GetClaim(ctx, id)
}

// CreateReversal creates a new reversal
func (store *SQLStore) CreateReversal(ctx context.Context, claimID uuid.UUID) (sqlc.Reversal, error) {
	return store.Queries.CreateReversal(ctx, claimID)
}

// CreatePharmacy creates a new pharmacy
func (store *SQLStore) CreatePharmacy(ctx context.Context, arg sqlc.CreatePharmacyParams) (sqlc.Pharmacy, error) {
	return store.Queries.CreatePharmacy(ctx, arg)
}

// GetPharmacy gets a pharmacy by NPI
func (store *SQLStore) GetPharmacy(ctx context.Context, npi string) (sqlc.Pharmacy, error) {
	return store.Queries.GetPharmacy(ctx, npi)
}

// CountPharmacies counts the total number of pharmacies
func (store *SQLStore) CountPharmacies(ctx context.Context) (int64, error) {
	return store.Queries.CountPharmacies(ctx)
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}
