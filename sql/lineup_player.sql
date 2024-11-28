-- name: FindLineupPlayerByLineupIDAndPositionNo :one
SELECT
    "lineup_player".*,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "player"
ON "lineup_player".player_id = "player".id
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".position_no = $2
LIMIT 1;

-- name: CreateLineupPlayer :one
INSERT INTO "lineup_player" (
    lineup_id,
    player_id,
    no,
    position_no,
    position
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: UpdateLineupPlayer :one
UPDATE lineup_player SET
    position_no = COALESCE(sqlc.narg('position_no'), position_no),
    position = COALESCE(sqlc.narg('position'), position)
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".player_id = $2
RETURNING *;

-- name: DeleteLineupPlayer :exec
DELETE FROM "lineup_player"
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".player_id = $2;

-- name: ListLineupPlayersByLineupID :many
SELECT
    "lineup_player".*,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "player"
ON "lineup_player".player_id = "player".id
WHERE "lineup_player".lineup_id = $1
ORDER BY "lineup_player".no ASC;