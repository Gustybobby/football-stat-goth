-- name: ListClubsOrderByNameAsc :many
SELECT *
FROM "club"
ORDER BY "club".name ASC;

-- name: FindClubByID :one
SELECT *
FROM "club"
WHERE "club".id = $1
LIMIT 1;

-- name: ClubAverageStatistics :one
SELECT
    AVG("lineup".possession) AS avg_possession,
    AVG("lineup".shots_on_target) AS avg_shots_on_target,
    AVG("lineup".shots) AS avg_shots,
    AVG("lineup".touches) AS avg_touches,
    AVG("lineup".passes) AS avg_passes,
    AVG("lineup".tackles) AS avg_tackles,
    AVG("lineup".clearances) AS avg_clearances,
     AVG("lineup".corners) AS avg_corners,
    AVG("lineup".offsides) AS avg_offsides,
    AVG("lineup".fouls_conceded) AS avg_fouls_conceded
FROM "lineup"
INNER JOIN "match"
ON "match".is_finished AND ("match".home_lineup_id = "lineup".id OR "match".away_lineup_id = "lineup".id)
WHERE "lineup".club_id = sqlc.arg(club_id)::text;