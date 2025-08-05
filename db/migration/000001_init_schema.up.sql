CREATE TABLE "claims" (
  "id" varchar PRIMARY KEY NOT NULL,
  "ndc" varchar NOT NULL,
  "quantity" int NOT NULL,
  "npi" varchar NOT NULL,
  "price" decimal(15,2),
  "timestamp" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "revert" (
  "id" varchar PRIMARY KEY NOT NULL,
  "claim_id" varchar NOT NULL,
  "timestamp" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "pharmacy" (
  "npi" varchar PRIMARY KEY NOT NULL,
  "chain" varchar NOT NULL
);

CREATE INDEX ON "claims" ("ndc");

CREATE INDEX ON "claims" ("npi");

CREATE INDEX ON "claims" ("ndc", "npi");

ALTER TABLE "revert" ADD FOREIGN KEY ("claim_id") REFERENCES "claims" ("id");

ALTER TABLE "pharmacy" ADD FOREIGN KEY ("npi") REFERENCES "claims" ("npi");
