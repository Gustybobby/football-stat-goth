-- name: CreateLineupPlayer :one
INSERT INTO "lineup_player" (
    lineup_id,
    player_id,
    position_no,
    position
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: ListLineupPlayersByLineupID :many
SELECT
    "lineup_player".*,
    "player".no,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "player"
ON "lineup_player".player_id = "player".id
WHERE "lineup_player".lineup_id = $1;