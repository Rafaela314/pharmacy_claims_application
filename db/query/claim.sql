INSERT INTO claims (
  id, ndc, quantity, npi, price     
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;
