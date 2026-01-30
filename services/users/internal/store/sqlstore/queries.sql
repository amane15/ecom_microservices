-- name: GetUser :one
SELECT
    id,
    email,
    first_name,
    last_name,
    ROLE,
    created_at,
    updated_at
FROM
    users
WHERE
    id = $1;

-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name)
    VALUES ($1, $2, $3)
RETURNING
    id, email, first_name, last_name, role, created_at, updated_at;

-- name: UpdateUser :one
UPDATE
    users
SET
    first_name = COALESCE(sqlc.narg (first_name), first_name),
    last_name = COALESCE(sqlc.narg (last_name), last_name)
WHERE
    id = $1
RETURNING
    id,
    email,
    first_name,
    last_name,
    role,
    created_at,
    updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

