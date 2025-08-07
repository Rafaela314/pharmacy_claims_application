-- name: CreateClaim :one
INSERT INTO claims (
  ndc, quantity, npi, price     
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetClaim :one
SELECT * FROM claims
WHERE id = $1 LIMIT 1;
