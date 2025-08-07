-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE pharmacies (
  npi VARCHAR PRIMARY KEY NOT NULL,
  chain VARCHAR NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE claims (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  ndc VARCHAR NOT NULL,
  quantity BIGINT NOT NULL,
  npi VARCHAR NOT NULL REFERENCES pharmacies(npi) ON DELETE CASCADE,
  price DOUBLE PRECISION NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE reversals (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX claims_ndc_idx ON claims (ndc);
CREATE INDEX claims_npi_idx ON claims (npi);
CREATE INDEX claims_ndc_npi_idx ON claims (ndc, npi);