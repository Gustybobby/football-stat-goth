-- name: ListMatchesWithClubsAndGoals :many
SELECT
    "match".*,
    "home_club".id AS home_club_id,
    "home_club".logo AS home_club_logo,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE (
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "match".home_lineup_id
        ) OR (
            "lineup_event"."event" = 'OWN_GOAL' AND
            "lineup_event".lineup_id = "match".away_lineup_id
        )
    ) AS home_goals,
    "away_club".id AS away_club_id,
    "away_club".logo AS away_club_logo,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE (
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "match".away_lineup_id
        ) OR (
            "lineup_event"."event" = 'OWN_GOAL' AND
            "lineup_event".lineup_id = "match".home_lineup_id
        )
    ) AS away_goals
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
    is_finished = sqlc.arg('is_finished')::bool AND
    CASE
        WHEN sqlc.arg('filter_club_id')::bool
        THEN
            "home_club".id = sqlc.arg('club_id')::TEXT OR
            "away_club".id = sqlc.arg('club_id')::TEXT
        ELSE true
    END AND
    CASE
        WHEN sqlc.arg('filter_week')::bool
        THEN "match".week = sqlc.arg('week')::INTEGER
        ELSE true
    END
ORDER BY
    CASE WHEN sqlc.arg('order')::TEXT = 'ASC' THEN "match".start_at END ASC,
    CASE WHEN sqlc.arg('order')::TEXT = 'DESC' THEN "match".start_at END DESC;

-- name: FindLatestFinishedMatchweek :one
SELECT
    "match".week
FROM "match"
WHERE "match".is_finished = true
ORDER BY "match".start_at DESC;

-- name: FindMatchByID :one
SELECT
    "match".*,
    "home_club".id AS home_club_id,
    "home_club".short_name AS home_club_name,
    "home_club".logo AS home_club_logo,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE (
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "match".home_lineup_id
        ) OR (
            "lineup_event"."event" = 'OWN_GOAL' AND
            "lineup_event".lineup_id = "match".away_lineup_id
        )
    ) AS home_goals,
    "home_lineup".possession AS home_possession,
    "home_lineup".shots_on_target AS home_shots_on_target,
    "home_lineup".shots AS home_shots,
    "home_lineup".touches AS home_touches,
    "home_lineup".passes AS home_passes,
    "home_lineup".tackles AS home_tackles,
    "home_lineup".clearances AS home_clearances,
    "home_lineup".corners AS home_corners,
    "home_lineup".offsides AS home_offsides,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'YELLOW' AND
            "lineup_event".lineup_id = "match".home_lineup_id
    ) AS home_yellow_cards,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'RED' AND
            "lineup_event".lineup_id = "match".home_lineup_id
    ) AS home_red_cards,
    "home_lineup".fouls_conceded AS home_fouls_conceded,
    "away_club".id AS away_club_id,
    "away_club".short_name AS away_club_name,
    "away_club".logo AS away_club_logo,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE (
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "match".away_lineup_id
        ) OR (
            "lineup_event"."event" = 'OWN_GOAL' AND
            "lineup_event".lineup_id = "match".home_lineup_id
        )
    ) AS away_goals,
    "away_lineup".possession AS away_possession,
    "away_lineup".shots_on_target AS away_shots_on_target,
    "away_lineup".shots AS away_shots,
    "away_lineup".touches AS away_touches,
    "away_lineup".passes AS away_passes,
    "away_lineup".tackles AS away_tackles,
    "away_lineup".clearances AS away_clearances,
    "away_lineup".corners AS away_corners,
    "away_lineup".offsides AS away_offsides,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'YELLOW' AND
            "lineup_event".lineup_id = "match".away_lineup_id
    ) AS away_yellow_cards,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE
            "lineup_event"."event" = 'RED' AND
            "lineup_event".lineup_id = "match".away_lineup_id
    ) AS away_red_cards,
    "away_lineup".fouls_conceded AS away_fouls_conceded
FROM "match"
INNER JOIN "lineup" as "home_lineup"
ON "match".home_lineup_id = "home_lineup".id
INNER JOIN "club" as "home_club"
ON "home_lineup".club_id = "home_club".id
INNER JOIN "lineup" as "away_lineup"
ON "match".away_lineup_id = "away_lineup".id
INNER JOIN "club" as "away_club"
ON "away_lineup".club_id = "away_club".id
WHERE "match".id = $1
LIMIT 1;

-- name: FindMatchIDFromLineupID :one
SELECT
    "match".id
FROM "match"
WHERE
    "match".home_lineup_id = sqlc.arg('lineup_id')::integer OR
    "match".away_lineup_id = sqlc.arg('lineup_id')::integer
LIMIT 1;

-- name: ListClubStandings :many
WITH "match_score" AS (
    SELECT
        "match".id,
        "home".club_id AS home_club_id,
        (
            SELECT COUNT(*)
            FROM "lineup_event"
            WHERE (
                "lineup_event"."event" = 'GOAL' AND
                "lineup_event".lineup_id = "match".home_lineup_id
            ) OR (
                "lineup_event"."event" = 'OWN_GOAL' AND
                "lineup_event".lineup_id = "match".away_lineup_id
            )
        ) AS home_goals,
        "away".club_id AS away_club_id,
        (
            SELECT COUNT(*)
            FROM "lineup_event"
            WHERE (
                "lineup_event"."event" = 'GOAL' AND
                "lineup_event".lineup_id = "match".away_lineup_id
            ) OR (
                "lineup_event"."event" = 'OWN_GOAL' AND
                "lineup_event".lineup_id = "match".home_lineup_id
            )
        ) AS away_goals,
        (
            SELECT COUNT(*)
            FROM "lineup_event"
            WHERE (
                "lineup_event"."event" = 'GOAL' AND
                "lineup_event".lineup_id = "match".home_lineup_id
            ) OR (
                "lineup_event"."event" = 'OWN_GOAL' AND
                "lineup_event".lineup_id = "match".away_lineup_id
            )
        ) - (
            SELECT COUNT(*)
            FROM "lineup_event"
            WHERE (
                "lineup_event"."event" = 'GOAL' AND
                "lineup_event".lineup_id = "match".away_lineup_id
            ) OR (
                "lineup_event"."event" = 'OWN_GOAL' AND
                "lineup_event".lineup_id = "match".home_lineup_id
            )
        ) AS goals_diff
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
    SUM("results".wins) AS won,
    SUM("results".draws) AS drawn,
    SUM("results".losses) AS lost,
    SUM("results".goals) AS gf,
    SUM("results".opp_goals) AS ga
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
ORDER BY
    SUM("results".wins) * 3 + SUM("results".draws) * 1 + SUM("results".losses) * 0 DESC,
    SUM("results".goals) - SUM("results".opp_goals) DESC;

-- name: ListPlayerMatchHistory :many
SELECT
    "match".*,
    "home_club".id AS home_club_id,
    "home_club".logo AS home_club_logo,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE (
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "match".home_lineup_id
        ) OR (
            "lineup_event"."event" = 'OWN_GOAL' AND
            "lineup_event".lineup_id = "match".away_lineup_id
        )
    ) AS home_goals,
    "away_club".id AS away_club_id,
    "away_club".logo AS away_club_logo,
    (
        SELECT COUNT(*)
        FROM "lineup_event"
        WHERE (
            "lineup_event"."event" = 'GOAL' AND
            "lineup_event".lineup_id = "match".away_lineup_id
        ) OR (
            "lineup_event"."event" = 'OWN_GOAL' AND
            "lineup_event".lineup_id = "match".home_lineup_id
        )
    ) AS away_goals,
    "lineup_player".position AS player_position,
    "lineup_player".no AS player_no
FROM "match"
INNER JOIN "lineup" as "home_lineup"
ON "match".home_lineup_id = "home_lineup".id
INNER JOIN "club" as "home_club"
ON "home_lineup".club_id = "home_club".id
INNER JOIN "lineup" as "away_lineup"
ON "match".away_lineup_id = "away_lineup".id
INNER JOIN "club" as "away_club"
ON "away_lineup".club_id = "away_club".id
INNER JOIN "lineup_player"
ON "lineup_player".player_id = sqlc.arg('player_id')::INTEGER AND (
    "lineup_player".lineup_id = "home_lineup".id OR
    "lineup_player".lineup_id = "away_lineup".id
)
WHERE is_finished = true
ORDER BY "match".start_at DESC;