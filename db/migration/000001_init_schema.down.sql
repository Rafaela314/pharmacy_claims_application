-- Drop indexes
DROP INDEX IF EXISTS claims_ndc_idx;
DROP INDEX IF EXISTS claims_npi_idx;
DROP INDEX IF EXISTS claims_ndc_npi_idx;

-- Drop tables in order to avoid FK conflicts
DROP TABLE IF EXISTS reversals;
DROP TABLE IF EXISTS claims;
DROP TABLE IF EXISTS pharmacies;

-- Optionally, drop the UUID extension if no longer needed
DROP EXTENSION IF EXISTS "uuid-ossp";