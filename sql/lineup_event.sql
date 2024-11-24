-- name: ListLineupEventsByMatchID :many
SELECT
    "lineup_event".*,
    "club_player1".no AS p1_no,
    "player1".firstname AS p1_firstname,
    "player1".lastname AS p1_lastname,
    "club_player2".no AS p2_no,
    "player2".firstname AS p2_firstname,
    "player2".lastname AS p2_lastname
FROM "lineup_event"
INNER JOIN "lineup"
ON "lineup_event".lineup_id = "lineup".id
INNER JOIN "match"
ON "lineup".id IN ("match".home_lineup_id, "match".away_lineup_id)
LEFT JOIN "club_player" AS "club_player1"
ON
    "lineup".club_id = "club_player1".club_id AND
    "lineup_event".player_id1 = "club_player1".player_id AND
    "match".season = "club_player1".season
LEFT JOIN "player" AS "player1"
ON "club_player1".player_id = "player1".id
LEFT JOIN "club_player" AS "club_player2"
ON
    "lineup".club_id = "club_player2".club_id AND
    "lineup_event".player_id2 = "club_player2".player_id AND
    "match".season = "club_player2".season
LEFT JOIN "player" AS "player2"
ON "club_player2".player_id = "player2".id
WHERE "match".id = $1
ORDER BY ("lineup_event".minutes + COALESCE("lineup_event".extra,0)) ASC;

    