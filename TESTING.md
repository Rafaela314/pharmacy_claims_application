# Testing Strategy

This document outlines the testing strategy for the pharmacy claims application to prevent dirty data and ensure test isolation.

## Overview

The application uses a comprehensive testing strategy with transaction-based testing to prevent dirty data in the database:

1. **Transaction-based testing** - Automatic rollback using database transactions
2. **Isolated test data** - Each test creates its own data within the transaction

## Testing Approaches

### Transaction-based Testing

Uses the `runTestWithTransaction` function to run tests within a database transaction that gets automatically rolled back.

```go
func TestCreateClaimWithTransaction(t *testing.T) {
    runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
        // Your test code here using txQueries instead of testQueries
        pharmacy, err := txQueries.CreatePharmacy(context.Background(), pharmacyArg)
        require.NoError(t, err)
        
        claim, err := txQueries.CreateClaim(context.Background(), claimArg)
        require.NoError(t, err)
        
        // Test assertions...
    })
    // Transaction is automatically rolled back here
}
```

**Benefits:**
- Fastest approach
- Automatic cleanup (rollback)
- Perfect isolation
- No manual cleanup required
- All database operations are automatically rolled back

## Test Data Creation

### Helper Functions

The test suite provides several helper functions for creating test data:

- `createRandomPharmacy(t)` - Creates a random pharmacy
- `createRandomClaim(t)` - Creates a random claim with a pharmacy
- `createRandomReversal(t)` - Creates a random reversal with a claim
- `createTestData(t)` - Creates a complete set of test data (pharmacy, claim, reversal)

### Example Usage

```go
func TestCompleteWorkflow(t *testing.T) {
    runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
        // Create test data
        pharmacy, claim, reversal := createTestData(t)
        
        // Test the complete workflow
        require.NotEmpty(t, pharmacy.NPI)
        require.NotEmpty(t, claim.ID)
        require.Equal(t, claim.ID, reversal.ClaimID)
    })
}
```

## Best Practices

### 1. Always Use Transactions

Never run tests without proper transaction isolation:

```go
// ❌ Bad - No transaction
func TestBadExample(t *testing.T) {
    claim := createRandomClaim(t)
    // Test code...
}

// ✅ Good - With transaction
func TestGoodExample(t *testing.T) {
    runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
        // Create test data within transaction
        claimArg := CreateClaimParams{...}
        claim, err := txQueries.CreateClaim(context.Background(), claimArg)
        // Test code...
    })
}
```

### 2. Use Descriptive Test Names

```go
// ❌ Bad
func TestCreate(t *testing.T) { ... }

// ✅ Good
func TestCreateClaimWithValidData(t *testing.T) { ... }
```

### 3. Test Edge Cases

```go
func TestCreateClaimWithMinimumValues(t *testing.T) {
    runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
        // Test with minimum valid values
        arg := CreateClaimParams{
            NDC:      "123456789",
            Price:    0.01,
            Quantity: 1,
            NPI:      "1234567890",
        }
        claim, err := txQueries.CreateClaim(context.Background(), arg)
        require.NoError(t, err)
        require.Equal(t, int64(1), claim.Quantity)
    })
}
```

### 4. Test Error Conditions

```go
func TestCreateClaimWithInvalidData(t *testing.T) {
    runTestWithTransaction(t, func(t *testing.T, txQueries *Queries) {
        // Test with invalid data
        arg := CreateClaimParams{
            NDC:      "", // Invalid: empty NDC
            Price:    -1, // Invalid: negative price
            Quantity: 0,  // Invalid: zero quantity
            NPI:      "1234567890",
        }
        _, err := txQueries.CreateClaim(context.Background(), arg)
        require.Error(t, err) // Should fail
    })
}
```

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Specific Test Package
```bash
go test ./db/sqlc
```

### Run Specific Test
```bash
go test -run TestCreateClaim ./db/sqlc
```

### Run Tests with Verbose Output
```bash
go test -v ./db/sqlc
```

### Run Tests with Coverage
```bash
go test -cover ./db/sqlc
```

## Database Setup for Testing

### Current Setup

The tests currently use the same database as development but with comprehensive transaction-based testing:

1. **Transaction-based testing** - Uses database transactions that are automatically rolled back
2. **Automatic cleanup** - All test data is automatically cleaned up via transaction rollback

### Optional: Separate Test Database

For even better isolation, you can set up a separate test database:

#### 1. Create a Separate Test Database

```sql
CREATE DATABASE pharmacy_claims_test;
```

#### 2. Update Test Configuration

Add a test database configuration to your `app.env`:

```env
DB_SOURCE_TEST=postgresql://root:password@localhost:5432/pharmacy_claims_test?sslmode=disable
```

#### 3. Run Migrations on Test Database

```bash
migrate -path db/migration -database "postgresql://root:password@localhost:5432/pharmacy_claims_test?sslmode=disable" up
```

#### 4. Update Test Configuration

Change the test setup to use the test database:

```go
// In db/sqlc/main_test.go
conn, err := pgxpool.New(context.Background(), config.DBSourceTest)
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:latest
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: pharmacy_claims_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - run: go test -v ./...
```

## Troubleshooting

### Common Issues

1. **Tests failing due to existing data**
   - Ensure cleanup functions are working
   - Check that tests are using the correct database

2. **Transaction rollback errors**
   - Check that all operations support transactions
   - Verify database connection is healthy

3. **Foreign key constraint errors**
   - Ensure test data is created in the correct order
   - Check that cleanup deletes in reverse dependency order

### Debug Mode

Run tests with debug output:

```bash
go test -v -debug ./db/sqlc
```

This will show detailed information about test execution and database operations. 