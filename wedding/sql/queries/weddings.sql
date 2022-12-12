-- name: CreateWedding :one
INSERT INTO weddings (
    id, names, wedding_date, user_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING status, created_at, updated_at;

-- name: GetWedding :one
SELECT * from weddings
WHERE id = $1;
