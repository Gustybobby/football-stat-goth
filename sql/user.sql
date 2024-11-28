-- name: FindUserByUsername :one
SELECT
    "user".username,
    "user".firstname,
    "user".lastname,
    "user".role
FROM "user"
WHERE "user".username = $1
LIMIT 1;

-- name: FindPasswordHashByUsername :one
SELECT
    "user".password_hash
FROM "user"
WHERE "user".username = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (
    username, password_hash, firstname, lastname
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdatePasswordByUsername :exec
UPDATE 
    "user" 
SET password_hash = $2
WHERE "user".username = $1;