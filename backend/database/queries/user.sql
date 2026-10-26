-- name: CreateUser :one
INSERT INTO users (
    email,
    password
)
VALUES ($1, $2)
RETURNING *;

-- name: FindUserByEmail :one
SELECT
    id,
    email,
    password
FROM users
WHERE email = $1;