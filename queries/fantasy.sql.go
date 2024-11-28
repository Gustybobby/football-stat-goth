// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: fantasy.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFantasyTeam = `-- name: CreateFantasyTeam :one
INSERT INTO "fantasy_team" (
    username,
    season,
    budget
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, username, season, budget
`

type CreateFantasyTeamParams struct {
	Username string
	Season   string
	Budget   int32
}

func (q *Queries) CreateFantasyTeam(ctx context.Context, arg CreateFantasyTeamParams) (FantasyTeam, error) {
	row := q.db.QueryRow(ctx, createFantasyTeam, arg.Username, arg.Season, arg.Budget)
	var i FantasyTeam
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Season,
		&i.Budget,
	)
	return i, err
}

type CreateFantasyTransactionParams struct {
	Cost            int32
	Type            FantasyTransactionType
	FantasyTeamID   int32
	FantasyPlayerID int32
}

const findFantasyTeamByUsernameSeason = `-- name: FindFantasyTeamByUsernameSeason :one
SELECT id, username, season, budget
FROM "fantasy_team"
WHERE
    "fantasy_team".username = $1 AND
    "fantasy_team".season = $2
`

type FindFantasyTeamByUsernameSeasonParams struct {
	Username string
	Season   string
}

func (q *Queries) FindFantasyTeamByUsernameSeason(ctx context.Context, arg FindFantasyTeamByUsernameSeasonParams) (FantasyTeam, error) {
	row := q.db.QueryRow(ctx, findFantasyTeamByUsernameSeason, arg.Username, arg.Season)
	var i FantasyTeam
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Season,
		&i.Budget,
	)
	return i, err
}

const listFantasyPlayers = `-- name: ListFantasyPlayers :many
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
                        )
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
    WHERE "match".season = $5::TEXT
    GROUP BY "player".id
), "player_ranked_total_stats" AS (
    SELECT
        player_total_stats.id, player_total_stats.total_goals, player_total_stats.total_assists, player_total_stats.total_yellow_cards, player_total_stats.total_red_cards, player_total_stats.total_own_goals, player_total_stats.appearances, player_total_stats.clean_sheets,
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
    "player".lastname,
    "player".position,
    "player".image,
    "club".id AS club_id,
    (
        SELECT
            CAST(
                ROUND(
                    $1::INTEGER + (
                        ($2::INTEGER - $1::INTEGER) /
                        ("rank_stats".rank_count_ratio - 1)
                    ) * ("player_ranked_total_stats".fantasy_rev_rank - 1)
                ) AS INTEGER
            )
        FROM "rank_stats"
        INNER JOIN "player_ranked_total_stats"
        ON "player_ranked_total_stats".id = "player".id
    ) AS cost
FROM "fantasy_player" 
INNER JOIN "player"
ON "fantasy_player".player_id = "player".id
INNER JOIN "club"
ON "fantasy_player".club_id = "club".id
WHERE
    CASE
        WHEN $3::bool
        THEN "fantasy_player".id = ANY($4::INTEGER[])
        ELSE true
    END
ORDER BY
    "player".position ASC,
    "player".lastname ASC
`

type ListFantasyPlayersParams struct {
	MinCost               int32
	AvgCost               int32
	FilterFantasyPlayerID bool
	FantasyPlayerIds      []int32
	Season                string
}

type ListFantasyPlayersRow struct {
	ID       int32
	PlayerID int32
	Lastname string
	Position PlayerPosition
	Image    pgtype.Text
	ClubID   string
	Cost     int32
}

func (q *Queries) ListFantasyPlayers(ctx context.Context, arg ListFantasyPlayersParams) ([]ListFantasyPlayersRow, error) {
	rows, err := q.db.Query(ctx, listFantasyPlayers,
		arg.MinCost,
		arg.AvgCost,
		arg.FilterFantasyPlayerID,
		arg.FantasyPlayerIds,
		arg.Season,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFantasyPlayersRow
	for rows.Next() {
		var i ListFantasyPlayersRow
		if err := rows.Scan(
			&i.ID,
			&i.PlayerID,
			&i.Lastname,
			&i.Position,
			&i.Image,
			&i.ClubID,
			&i.Cost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFantasyTeamPlayersByUsernameSeason = `-- name: ListFantasyTeamPlayersByUsernameSeason :many
SELECT
    fantasy_team_player.fantasy_team_id, fantasy_team_player.fantasy_player_id,
    "fantasy_team".budget
FROM "fantasy_team_player"
INNER JOIN "fantasy_team"
ON "fantasy_team_player".fantasy_team_id = "fantasy_team".id
WHERE
    "fantasy_team".username = $1 AND
    "fantasy_team".season = $2
`

type ListFantasyTeamPlayersByUsernameSeasonParams struct {
	Username string
	Season   string
}

type ListFantasyTeamPlayersByUsernameSeasonRow struct {
	FantasyTeamID   int32
	FantasyPlayerID int32
	Budget          int32
}

func (q *Queries) ListFantasyTeamPlayersByUsernameSeason(ctx context.Context, arg ListFantasyTeamPlayersByUsernameSeasonParams) ([]ListFantasyTeamPlayersByUsernameSeasonRow, error) {
	rows, err := q.db.Query(ctx, listFantasyTeamPlayersByUsernameSeason, arg.Username, arg.Season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFantasyTeamPlayersByUsernameSeasonRow
	for rows.Next() {
		var i ListFantasyTeamPlayersByUsernameSeasonRow
		if err := rows.Scan(&i.FantasyTeamID, &i.FantasyPlayerID, &i.Budget); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
