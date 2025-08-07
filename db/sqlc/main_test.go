package db

import (
	"context"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pharmacy_claims_application/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../../pharmacy_claims_application")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Connect to test database
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
