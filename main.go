package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pharmacy_claims_application/db"
	"github.com/pharmacy_claims_application/seeder"
	"github.com/pharmacy_claims_application/server"
	"github.com/pharmacy_claims_application/util"
)

func main() {
	// Load configuration
	config, err := util.LoadConfig("")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Connect to database
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	// Create database store
	store := db.NewStore(conn)

	// Seed database with pharmacy data if empty
	if err := seeder.SeedPharmacies(store, "data"); err != nil {
		log.Printf("Warning: failed to seed pharmacies: %v", err)
	}

	// Create and start server
	server := server.NewServer(store)

	// Start server in a goroutine
	go func() {
		if err := server.Start(config); err != nil {
			log.Fatal("cannot start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
