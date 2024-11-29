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

-- name: ListUsers :many
SELECT
    "user".username,
    "user".firstname,
    "user".lastname,
    "user".role
FROM "user"
ORDER BY
    "user".username ASC,
    "user".role ASC;


-- name: CreateUser :one
INSERT INTO "user" (
    username, password_hash, firstname, lastname
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE "user" SET
    firstname = COALESCE(sqlc.narg('firstname')::TEXT, firstname),
    lastname = COALESCE(sqlc.narg('lastname')::TEXT, lastname)
WHERE "user".username = $1;

-- name: UpdatePasswordByUsername :exec
UPDATE 
    "user" 
SET password_hash = $2
WHERE "user".username = $1;