-- name: ListFantasyPlayers :many
WITH "player_total_stats" AS (
    SELECT
        "player".id,
        COUNT(
            CASE
                WHEN
                    "lineup_event".event = 'GOAL' AND
                    "lineup_event".player_id1 = "player".id
                THEN 1 ELSE NULL
            END
        ) AS total_goals,
        COUNT(
            CASE
                WHEN
                    "lineup_event".event = 'GOAL' AND
                    "lineup_event".player_id2 = "player".id
                THEN 1 ELSE NULL
            END
        ) AS total_assists,
        COUNT(
            CASE
                WHEN
                    "lineup_event".event = 'YELLOW' AND
                    "lineup_event".player_id1 = "player".id
                THEN 1 ELSE NULL
            END
        ) AS total_yellow_cards,
        COUNT(
            CASE
                WHEN
                    "lineup_event".event = 'RED' AND
                    "lineup_event".player_id1 = "player".id
                THEN 1 ELSE NULL
            END
        ) AS total_red_cards,
        COUNT(
            CASE
                WHEN
                    "lineup_event".event = 'OWN_GOAL' AND
                    "lineup_event".player_id1 = "player".id
                THEN 1 ELSE NULL
            END
        ) AS total_own_goals,
        COUNT(DISTINCT "lineup_player".lineup_id) AS appearances,
        CASE
            WHEN "player".position = 'GK' THEN (
                SELECT
                    COUNT(
                        CASE
                            WHEN EXISTS (
                                SELECT 1
                                FROM "lineup_player"
                                WHERE
                                    "lineup_player".player_id = "player".id AND
                                    "lineup_player".lineup_id = "match".home_lineup_id
                            ) THEN
                                CASE
                                    WHEN EXISTS (
                                        SELECT 1
                                        FROM "lineup_event"
                                        WHERE (
                                            "lineup_event".event = 'GOAL' AND
                                            "lineup_event".lineup_id = "match".away_lineup_id
                                        ) OR (
                                            "lineup_event".event = 'OWN_GOAL' AND
                                            "lineup_event".lineup_id = "match".home_lineup_id
                                        )
                                    ) THEN NULL ELSE 1
                                END
                            ELSE
                                CASE
                                    WHEN EXISTS (
                                        SELECT 1
                                        FROM "lineup_event"
                                        WHERE (
                                            "lineup_event".event = 'GOAL' AND
                                            "lineup_event".lineup_id = "match".home_lineup_id
                                        ) OR (
                                            "lineup_event".event = 'OWN_GOAL' AND
                                            "lineup_event".lineup_id = "match".away_lineup_id
                                        )
                                    ) THEN NULL ELSE 1
                                END
                        END
                    )
                FROM "match"
                WHERE EXISTS (
                    SELECT 1
                    FROM "lineup_player"
                    WHERE
                        "lineup_player".player_id = "player".id AND (
                            "lineup_player".lineup_id = "match".home_lineup_id OR
                            "lineup_player".lineup_id = "match".away_lineup_id
                        ) AND "lineup_player".position = 'GK'
                    )
            )
            ELSE 0
        END AS clean_sheets
    FROM "player"
    LEFT JOIN "lineup_player"
    ON "player".id = "lineup_player".player_id
    LEFT JOIN "lineup_event"
    ON
        "lineup_player".lineup_id = "lineup_event".lineup_id AND (
            "lineup_player".player_id = "lineup_event".player_id1 OR
            "lineup_player".player_id = "lineup_event".player_id2
        )
    LEFT JOIN "match"
    ON
        "lineup_player".lineup_id = "match".home_lineup_id OR
        "lineup_player".lineup_id = "match".away_lineup_id
    WHERE "match".season = sqlc.arg('season')::TEXT
    GROUP BY "player".id
), "player_ranked_total_stats" AS (
    SELECT
        "player_total_stats".*,
        RANK() OVER (
            ORDER BY "player_total_stats".total_goals DESC
        ) AS goals_rank,
        RANK() OVER (
            ORDER BY "player_total_stats".total_assists DESC
        ) AS assists_rank,
        RANK() OVER (
            ORDER BY (
                "player_total_stats".total_goals +
                "player_total_stats".clean_sheets +
                "player_total_stats".total_assists * 0.75 
            ) ASC
        ) AS fantasy_rev_rank
    FROM "player_total_stats"
), "rank_stats" AS (
    SELECT
        SUM("player_ranked_total_stats".fantasy_rev_rank)/COUNT(*) AS rank_count_ratio,
        CAST(MAX("player_ranked_total_stats".fantasy_rev_rank) AS INTEGER) AS max_rank
    FROM "player_ranked_total_stats"
)
SELECT
    "fantasy_player".id,
    "player".id AS player_id,
    "player".firstname,
    "player".lastname,
    "player".position,
    "player".image,
    "club".id AS club_id,
    CAST(
        COALESCE(
            (
                SELECT
                    CAST(
                        ROUND(
                            sqlc.arg('min_cost')::INTEGER + (
                                (sqlc.arg('avg_cost')::INTEGER - sqlc.arg('min_cost')::INTEGER) /
                                ("rank_stats".rank_count_ratio - 1)
                            ) * ("player_ranked_total_stats".fantasy_rev_rank - 1)
                        ) AS INTEGER
                    )
                FROM "rank_stats"
                INNER JOIN "player_ranked_total_stats"
                ON "player_ranked_total_stats".id = "player".id
            ),
            sqlc.arg('min_cost')::INTEGER
        ) AS INTEGER
    ) AS cost
FROM "fantasy_player" 
INNER JOIN "player"
ON "fantasy_player".player_id = "player".id
INNER JOIN "club"
ON "fantasy_player".club_id = "club".id
WHERE
    CASE
        WHEN sqlc.arg('filter_fantasy_player_id')::bool
        THEN "fantasy_player".id = ANY(sqlc.arg('fantasy_player_ids')::INTEGER[])
        ELSE true
    END
ORDER BY
    "player".position ASC,
    "player".lastname ASC;

-- name: FindFantasyTeamByUsernameSeason :one
SELECT *
FROM "fantasy_team"
WHERE
    "fantasy_team".username = $1 AND
    "fantasy_team".season = $2
LIMIT 1;

-- name: ListFantasyTeamPlayersByFantasyTeamID :many
SELECT *
FROM "fantasy_team_player"
WHERE "fantasy_team_player".fantasy_team_id = $1;

-- name: CountFantasyTeamPlayersByFantasyTeamID :one
SELECT
    COUNT(CASE WHEN "player".position = 'GK' THEN 1 ELSE NULL END) AS GK_count,
    COUNT(CASE WHEN "player".position = 'DEF' THEN 1 ELSE NULL END) AS DEF_count,
    COUNT(CASE WHEN "player".position = 'MFD' THEN 1 ELSE NULL END) AS MFD_count,
    COUNT(CASE WHEN "player".position = 'FWD' THEN 1 ELSE NULL END) AS FWD_count
FROM "fantasy_team_player"
INNER JOIN "fantasy_player"
ON "fantasy_team_player".fantasy_player_id = "fantasy_player".id
INNER JOIN "player"
ON "fantasy_player".player_id = "player".id
WHERE "fantasy_team_player".fantasy_team_id = $1;

-- name: FindLastestTransaction :one
SELECT *
FROM "fantasy_transaction"
WHERE
    "fantasy_transaction".fantasy_team_id = $1 AND
    "fantasy_transaction".fantasy_player_id = $2
ORDER BY "fantasy_transaction".created_at DESC
LIMIT 1;


-- name: CreateFantasyTeam :one
INSERT INTO "fantasy_team" (
    username,
    season,
    budget
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: CreateFantasyTransaction :copyfrom
INSERT INTO "fantasy_transaction" (
    cost,
    type,
    fantasy_team_id,
    fantasy_player_id
) VALUES (
    $1,
    $2,
    $3,
    $4
);