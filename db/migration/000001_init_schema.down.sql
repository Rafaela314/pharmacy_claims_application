-- Drop foreign keys first
ALTER TABLE "revert" DROP CONSTRAINT IF EXISTS revert_claim_id_fkey;
ALTER TABLE "pharmacy" DROP CONSTRAINT IF EXISTS pharmacy_npi_fkey;

-- Drop indexes
DROP INDEX IF EXISTS claims_ndc_idx;
DROP INDEX IF EXISTS claims_npi_idx;
DROP INDEX IF EXISTS claims_ndc_npi_idx;

-- Drop tables
DROP TABLE IF EXISTS "revert";
DROP TABLE IF EXISTS "pharmacy";
DROP TABLE IF EXISTS "claims";