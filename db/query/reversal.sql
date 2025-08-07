-- name: CreateReversal :one
INSERT INTO reversals (
  claim_id  
) VALUES (
  $1
)
RETURNING *;

-- name: GetReversalByClaimID :one
SELECT * FROM reversals
WHERE claim_id = $1 LIMIT 1;

-- name: DeleteReversal :exec
DELETE FROM reversals
WHERE id = $1;