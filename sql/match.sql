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
    THEN "home_club".id = sqlc.arg(club_id)::text OR "away_club".id = sqlc.arg(club_id)::text
    ELSE true
    END
ORDER BY
    CASE WHEN sqlc.arg('order')::text = 'ASC' THEN "match".start_at END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN "match".start_at END DESC;

-- name: ListClubStandings :many
WITH "match_score" AS (
    SELECT
        "match".id,
        "home".club_id AS home_club_id,
        "home".goals AS home_goals,
        "away".club_id AS away_club_id,
        "away".goals AS away_goals,
        "home".goals - "away".goals AS goals_diff
    FROM "match"
    INNER JOIN "lineup" AS "home"
    ON "match".home_lineup_id = "home".id
    INNER JOIN "lineup" AS "away"
    ON "match".away_lineup_id = "away".id
    WHERE "match".is_finished = true
)
SELECT
    "club".id,
    "club".Name,
    "club".Logo,
    CAST(SUM("results".wins) AS INTEGER) AS won,
    CAST(SUM("results".draws) AS INTEGER) AS drawn,
    CAST(SUM("results".losses) AS INTEGER) AS lost,
    CAST(SUM("results".goals) AS INTEGER) AS gf,
    CAST(SUM("results".opp_goals) AS INTEGER) AS ga
FROM (
    SELECT
        "match_score".home_club_id AS club_id,
        SUM("match_score".home_goals) AS goals,
        SUM("match_score".away_goals) AS opp_goals,
        SUM(CASE WHEN "match_score".goals_diff > 0 THEN 1 ELSE 0 END) AS wins,
        SUM(CASE WHEN "match_score".goals_diff = 0 THEN 1 ELSE 0 END) AS draws,
        SUM(CASE WHEN "match_score".goals_diff < 0 THEN 1 ELSE 0 END ) AS losses
    FROM "match_score"
    GROUP BY home_club_id
    UNION ALL
    SELECT
        "match_score".away_club_id AS club_id,
        SUM("match_score".away_goals) AS goals,
        SUM("match_score".home_goals) AS opp_goals,
        SUM(CASE WHEN "match_score".goals_diff < 0 THEN 1 ELSE 0 END) AS wins,
        SUM(CASE WHEN "match_score".goals_diff = 0 THEN 1 ELSE 0 END) AS draws,
        SUM(CASE WHEN "match_score".goals_diff > 0 THEN 1 ELSE 0 END ) AS losses
    FROM "match_score"
    GROUP BY away_club_id
) AS "results"
INNER JOIN "club"
ON "results".club_id = club.id
GROUP BY "club".id
ORDER BY CAST(
    CAST(SUM("results".wins) AS INTEGER) * 3 +
    CAST(SUM("results".draws) AS INTEGER) * 1 +
    CAST(SUM("results".losses) AS INTEGER) * 0
    AS INTEGER
) DESC;