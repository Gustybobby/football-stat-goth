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
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "lineup_player".lineup_id
    ) AS goals,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'YELLOW' AND
            "lineup_event".lineup_id = "lineup_player".lineup_id
    ) AS yellow_cards,
     (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'RED' AND
            "lineup_event".lineup_id = "lineup_player".lineup_id
    ) AS red_cards,
    "player".no,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "player"
ON "lineup_player".player_id = "player".id
WHERE "lineup_player".lineup_id = $1;

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
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "lineup_player".lineup_id
    ) AS goals,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'YELLOW' AND
            "lineup_event".lineup_id = "lineup_player".lineup_id
    ) AS yellow_cards,
     (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'RED' AND
            "lineup_event".lineup_id = "lineup_player".lineup_id
    ) AS red_cards,
    "player".no,
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