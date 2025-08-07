package seeder

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pharmacy_claims_application/db"
	sqlc "github.com/pharmacy_claims_application/db/sqlc"
)

// PharmacyData represents a pharmacy record from CSV
type PharmacyData struct {
	Chain string `csv:"chain"`
	NPI   string `csv:"npi"`
}

// SeedPharmacies seeds the database with pharmacy data from CSV files
func SeedPharmacies(store db.Store, dataDir string) error {
	// Check if pharmacies table is empty
	count, err := store.CountPharmacies(context.Background())
	if err != nil {
		return fmt.Errorf("failed to check pharmacy count: %w", err)
	}

	if count > 0 {
		log.Println("Pharmacies table is not empty, skipping seed")
		return nil
	}

	// Find CSV files in the data/pharmacies directory
	csvFiles, err := findCSVFiles(filepath.Join(dataDir, "pharmacies"))
	if err != nil {
		return fmt.Errorf("failed to find CSV files: %w", err)
	}

	if len(csvFiles) == 0 {
		return fmt.Errorf("no CSV files found in %s", filepath.Join(dataDir, "pharmacies"))
	}

	// Process each CSV file
	for _, csvFile := range csvFiles {
		if err := processPharmacyCSV(store, csvFile); err != nil {
			return fmt.Errorf("failed to process %s: %w", csvFile, err)
		}
	}

	log.Printf("Successfully seeded %d pharmacies from CSV files", count)
	return nil
}

// findCSVFiles finds all CSV files in the given directory
func findCSVFiles(dir string) ([]string, error) {
	var csvFiles []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".csv") {
			csvFiles = append(csvFiles, filepath.Join(dir, entry.Name()))
		}
	}

	return csvFiles, nil
}

// processPharmacyCSV processes a single CSV file and inserts pharmacy data
func processPharmacyCSV(store db.Store, csvFile string) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or missing data")
	}

	// Skip header row
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		if len(record) < 2 {
			log.Printf("Skipping invalid record at line %d: %v", i+1, record)
			continue
		}

		pharmacy := PharmacyData{
			Chain: strings.TrimSpace(record[0]),
			NPI:   strings.TrimSpace(record[1]),
		}

		// Validate data
		if pharmacy.Chain == "" || pharmacy.NPI == "" {
			log.Printf("Skipping invalid pharmacy data at line %d: chain=%s, npi=%s", i+1, pharmacy.Chain, pharmacy.NPI)
			continue
		}

		// Insert into database
		arg := sqlc.CreatePharmacyParams{
			Chain: pharmacy.Chain,
			NPI:   pharmacy.NPI,
		}

		_, err := store.CreatePharmacy(context.Background(), arg)
		if err != nil {
			log.Printf("Failed to insert pharmacy %s: %v", pharmacy.NPI, err)
			continue // Continue with other records even if one fails
		}

		log.Printf("Inserted pharmacy: %s (NPI: %s)", pharmacy.Chain, pharmacy.NPI)
	}

	return nil
}
