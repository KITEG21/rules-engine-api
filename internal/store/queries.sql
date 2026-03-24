-- name: CreateRule :one
INSERT INTO rules (name, description, definition)
VALUES ($1, $2, $3)
RETURNING id, name, description, definition, is_active, created_at, updated_at;

-- name: GetRule :one
SELECT id, name, description, definition, is_active, created_at, updated_at
FROM rules
WHERE id = $1;

-- name: ListActiveRules :many
SELECT id, name, description, definition, is_active, created_at, updated_at
FROM rules
WHERE is_active = true
ORDER BY name;

-- name: UpdateRule :one
UPDATE rules
SET name = $2, description = $3, definition = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, definition, is_active, created_at, updated_at;

-- name: DeleteRule :exec
DELETE FROM rules WHERE id = $1;
