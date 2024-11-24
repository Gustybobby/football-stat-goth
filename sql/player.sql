-- name: ListPlayersOrderByNameAsc :many
SELECT *
FROM "player"
ORDER BY "player".club_id ASC;

-- name: FindPlayerIDByClubAndNo :one
SELECT
    "player".id
FROM "player"
WHERE
    "player".club_id = $1 AND
    "player".no = $2;

-- name: CreatePlayer :one
INSERT INTO "player" (
    club_id,
    no,
    firstname,
    lastname,
    dob,
    height,
    nationality,
    position,
    image
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;