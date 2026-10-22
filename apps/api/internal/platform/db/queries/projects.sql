-- name: GetProjectByPublicKey :one
SELECT * FROM projects WHERE public_key = $1;

-- name: CreateProject :one
INSERT INTO projects (owner_id, name, public_key) VALUES ($1, $2, $3) RETURNING *;