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
    "club_player".no,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "lineup"
ON "lineup_player".lineup_id = "lineup".id
INNER JOIN "match"
ON
    "lineup".id = "match".home_lineup_id OR
    "lineup".id = "match".away_lineup_id
INNER JOIN "club_player"
ON
    "lineup".club_id = "club_player".club_id AND
    "lineup_player".player_id = "club_player".player_id AND
    "match".season = "club_player".season
INNER JOIN "player"
ON "club_player".player_id = "player".id
WHERE "lineup_player".lineup_id = $1
ORDER BY "club_player".no ASC;

-- name: UpdateLineupPlayer :one
UPDATE lineup_player SET
    position_no = COALESCE(sqlc.narg('position_no'), position_no),
    position = COALESCE(sqlc.narg('position'), position)
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".player_id = $2
RETURNING *;

-- name: FindLineupPlayerByLineupIDAndPositionNo :one
SELECT
    "lineup_player".*,
    "club_player".no,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "lineup"
ON "lineup_player".lineup_id = "lineup".id
INNER JOIN "match"
ON
    "lineup".id = "match".home_lineup_id OR
    "lineup".id = "match".away_lineup_id
INNER JOIN "club_player"
ON
    "lineup".club_id = "club_player".club_id AND
    "lineup_player".player_id = "club_player".player_id AND
    "match".season = "club_player".season
INNER JOIN "player"
ON "club_player".player_id = "player".id
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".position_no = $2
LIMIT 1;