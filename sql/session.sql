-- name: FindSessionByToken :one
SELECT *
FROM "session"
WHERE "session".token = $1
LIMIT 1;

-- name: CreateSession :one
INSERT INTO "session" (
    token, username, expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateSessionExpiresAt :one
UPDATE "session"
SET expires_at = sqlc.arg(expires_at)::timestamp
WHERE "session".token = sqlc.arg(token)::text
RETURNING *; 