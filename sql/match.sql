-- name: ListMatchesWithClubsAndGoals :many
SELECT
    "match".*,
    "home_club".id AS home_club_id,
    "home_club".logo AS home_club_logo,
    "home_lineup".goals AS home_goals,
    "away_club".id AS away_club_id,
    "away_club".logo AS away_club_logo,
    "away_lineup".goals AS away_goals
FROM "match"
INNER JOIN "lineup" as "home_lineup"
ON "match".home_lineup_id = "home_lineup".id
INNER JOIN "club" as "home_club"
ON "home_lineup".club_id = "home_club".id
INNER JOIN "lineup" as "away_lineup"
ON "match".away_lineup_id = "away_lineup".id
INNER JOIN "club" as "away_club"
ON "away_lineup".club_id = "away_club".id
WHERE
    is_finished = sqlc.arg(is_finished)::bool AND
    CASE WHEN sqlc.arg(filter_club_id)::bool
    THEN "home_club".id = sqlc.arg(club_id)::text
    ELSE true
    END
ORDER BY
    CASE WHEN sqlc.arg('order')::text = 'ASC' THEN "match".start_at END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN "match".start_at END DESC;