-- name: FindPlayerByID :one
SELECT *
FROM "player"
WHERE "player".id = $1
LIMIT 1;

-- name: CreatePlayer :one
INSERT INTO "player" (
    firstname,
    lastname,
    dob,
    height,
    nationality,
    position,
    image
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: UpdatePlayerByID :one
UPDATE "player" SET
    firstname = $1,
    lastname = $2,
    dob = $3,
    height = $4,
    nationality = $5,
    position = $6,
    image = $7
WHERE "player".id = $8
RETURNING *;

-- name: DeletePlayerByID :exec
DELETE FROM "player"
WHERE "player".id = $1;

-- name: ListPlayerLikeFullname :many
SELECT *
FROM "player"
WHERE
    LOWER(CONCAT("player".firstname,' ',"player".lastname))
    LIKE LOWER(sqlc.arg('fullname_like')::TEXT)
ORDER BY "player".id ASC
OFFSET sqlc.arg('offset')::INTEGER
LIMIT sqlc.arg('limit')::INTEGER;

-- name: ListPlayerSeasonPerformance :many
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
    HAVING
        CASE
            WHEN sqlc.arg('filter_club_id')::bool
            THEN EXISTS (
                SELECT 1
                FROM "club_player"
                WHERE
                    "club_player".player_id = "player".id AND
                    "club_player".club_id = sqlc.arg('club_id')::TEXT AND
                    "club_player".season = sqlc.arg('season')::TEXT
            )
            ELSE true
        END
), "player_ranked_total_stats" AS (
    SELECT
        "player_total_stats".*,
        sqlc.arg('season')::TEXT AS season,
        RANK() OVER (
            ORDER BY "player_total_stats".total_goals DESC
        ) AS goals_rank,
        RANK() OVER (
            ORDER BY "player_total_stats".total_assists DESC
        ) AS assists_rank,
        RANK() OVER (
            ORDER BY "player_total_stats".clean_sheets DESC
        ) AS clean_sheets_rank,
        RANK() OVER (
            ORDER BY (
                "player_total_stats".total_goals +
                "player_total_stats".clean_sheets +
                "player_total_stats".total_assists * 0.75 
            ) DESC
        ) AS fantasy_rank
    FROM "player_total_stats"
)
SELECT *
FROM "player_ranked_total_stats"
WHERE
    CASE
        WHEN sqlc.arg('filter_player_id')::bool
        THEN "player_ranked_total_stats".id = sqlc.arg('player_id')::INTEGER
        ELSE true
    END
ORDER BY
    CASE
        WHEN sqlc.arg('order_by')::TEXT = 'GOAL'
        THEN "player_ranked_total_stats".goals_rank
        WHEN sqlc.arg('order_by')::TEXT = 'ASSIST'
        THEN "player_ranked_total_stats".assists_rank
        WHEN sqlc.arg('order_by')::TEXT = 'CLEANSHEET'
        THEN "player_ranked_total_stats".clean_sheets_rank
        WHEN sqlc.arg('order_by')::TEXT = 'FANTASY'
        THEN "player_ranked_total_stats".fantasy_rank
        ELSE NULL
    END ASC
LIMIT sqlc.narg('limit')::INTEGER;