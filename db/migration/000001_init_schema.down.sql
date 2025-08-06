-- Drop foreign key constraints
ALTER TABLE "reversals" DROP CONSTRAINT IF EXISTS reversals_claim_id_fkey;
ALTER TABLE "pharmacies" DROP CONSTRAINT IF EXISTS pharmacies_claim_id_fkey;

-- Drop indexes
DROP INDEX IF EXISTS claims_ndc_idx;
DROP INDEX IF EXISTS claims_npi_idx;
DROP INDEX IF EXISTS claims_ndc_npi_idx;

-- Drop tables
DROP TABLE IF EXISTS "reversals";
DROP TABLE IF EXISTS "pharmacies";
DROP TABLE IF EXISTS "claims";