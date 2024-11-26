-- name: FindPlayerByID :one
SELECT *
FROM "player"
WHERE "player".id = $1
LIMIT 1;

-- name: FindPlayerIDByClubNoSeason :one
SELECT
    "club_player".player_id
FROM "club_player"
WHERE
    "club_player".club_id = $1 AND
    "club_player".no = $2 AND
    "club_player".season = $3;

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

-- name: ListPlayerSeasonPerformance :many
WITH "total_stats" AS (
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
        (
            SELECT COUNT(*)
            FROM "match"
            WHERE EXISTS (
                SELECT 1
                FROM "lineup_player"
                WHERE
                    "lineup_player".player_id = "player".id AND (
                        "match".home_lineup_id = "lineup_player".lineup_id OR
                        "match".away_lineup_id = "lineup_player".lineup_id
                    )
            )
        ) AS appearances
    FROM "player"
    LEFT JOIN "lineup_event"
    ON
        "player".id = "lineup_event".player_id1 OR
        "player".id = "lineup_event".player_id2
    LEFT JOIN "match"
    ON
        "match".season = sqlc.arg('season')::TEXT AND (
            "lineup_event".lineup_id = "match".home_lineup_id OR
            "lineup_event".lineup_id = "match".away_lineup_id
        )
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
), "total_rank_stats" AS (
    SELECT
        "total_stats".*,
        sqlc.arg('season')::TEXT AS season,
        RANK() OVER (
            ORDER BY "total_stats".total_goals DESC
        ) AS goals_rank,
        RANK() OVER (
            ORDER BY "total_stats".total_assists DESC
        ) AS assists_rank
    FROM "total_stats"
)
SELECT *
FROM "total_rank_stats"
WHERE
    CASE
        WHEN sqlc.arg('filter_player_id')::bool
        THEN "total_rank_stats".id = sqlc.arg('player_id')::INTEGER
        ELSE true
    END
LIMIT sqlc.arg('limit')::INTEGER;