-- name: FindPlayerIDByClubNoSeason :one
SELECT
    "club_player".player_id
FROM "club_player"
WHERE
    "club_player".club_id = $1 AND
    "club_player".no = $2 AND
    "club_player".season = $3;

-- name: CreatePlayer :one
INSERT INTO "player" (
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
    $7
)
RETURNING *;