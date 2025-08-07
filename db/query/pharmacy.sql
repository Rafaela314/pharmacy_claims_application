-- name: CreatePharmacy :one
INSERT INTO pharmacies (
  npi, chain
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetPharmacy :one
SELECT * FROM pharmacies
WHERE npi = $1 LIMIT 1;

-- name: CountPharmacies :one
SELECT COUNT(*) FROM pharmacies;