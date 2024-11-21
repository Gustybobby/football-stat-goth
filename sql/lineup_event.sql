-- name: ListLineupEventsByMatchID :many
SELECT
    "lineup_event".*,
    "player1".no AS p1_no,
    "player1".firstname AS p1_firstname,
    "player1".lastname AS p1_lastname,
    "player2".no AS p2_no,
    "player2".firstname AS p2_firstname,
    "player2".lastname AS p2_lastname
FROM "lineup_event"
INNER JOIN "match"
ON
    "lineup_event".lineup_id = "match".home_lineup_id OR
    "lineup_event".lineup_id = "match".away_lineup_id
LEFT JOIN "player" AS "player1"
ON "lineup_event".player_id1 = "player1".id
LEFT JOIN "player" AS "player2"
ON "lineup_event".player_id2 = "player2".id
WHERE "match".id = $1
ORDER BY ("lineup_event".minutes + COALESCE("lineup_event".extra,0)) ASC;

    