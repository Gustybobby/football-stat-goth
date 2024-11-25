-- name: ListClubPlayerByPlayerID :many
SELECT
    "club_player".*,
    "club".short_name AS club_short_name,
    "club".logo AS club_logo
FROM "club_player"
INNER JOIN "club"
ON "club_player".club_id = "club".id
WHERE "club_player".player_id = $1
ORDER BY CAST(
    SPLIT_PART("club_player".season,'/',1) AS INTEGER
) DESC;