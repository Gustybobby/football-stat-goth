// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: player.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPlayer = `-- name: CreatePlayer :one
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
RETURNING id, firstname, lastname, dob, height, nationality, position, image
`

type CreatePlayerParams struct {
	Firstname   string
	Lastname    string
	Dob         pgtype.Timestamp
	Height      int16
	Nationality string
	Position    PlayerPosition
	Image       pgtype.Text
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (Player, error) {
	row := q.db.QueryRow(ctx, createPlayer,
		arg.Firstname,
		arg.Lastname,
		arg.Dob,
		arg.Height,
		arg.Nationality,
		arg.Position,
		arg.Image,
	)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Dob,
		&i.Height,
		&i.Nationality,
		&i.Position,
		&i.Image,
	)
	return i, err
}

const deletePlayerByID = `-- name: DeletePlayerByID :exec
DELETE FROM "player"
WHERE "player".id = $1
`

func (q *Queries) DeletePlayerByID(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deletePlayerByID, id)
	return err
}

const findPlayerByID = `-- name: FindPlayerByID :one
SELECT id, firstname, lastname, dob, height, nationality, position, image
FROM "player"
WHERE "player".id = $1
LIMIT 1
`

func (q *Queries) FindPlayerByID(ctx context.Context, id int32) (Player, error) {
	row := q.db.QueryRow(ctx, findPlayerByID, id)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Dob,
		&i.Height,
		&i.Nationality,
		&i.Position,
		&i.Image,
	)
	return i, err
}

const listPlayerLikeFullname = `-- name: ListPlayerLikeFullname :many
SELECT id, firstname, lastname, dob, height, nationality, position, image
FROM "player"
WHERE
    LOWER(CONCAT("player".firstname,' ',"player".lastname))
    LIKE LOWER($1::TEXT)
ORDER BY "player".id ASC
OFFSET $2::INTEGER
LIMIT $3::INTEGER
`

type ListPlayerLikeFullnameParams struct {
	FullnameLike string
	Offset       int32
	Limit        int32
}

func (q *Queries) ListPlayerLikeFullname(ctx context.Context, arg ListPlayerLikeFullnameParams) ([]Player, error) {
	rows, err := q.db.Query(ctx, listPlayerLikeFullname, arg.FullnameLike, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Player
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.ID,
			&i.Firstname,
			&i.Lastname,
			&i.Dob,
			&i.Height,
			&i.Nationality,
			&i.Position,
			&i.Image,
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

const listPlayerSeasonPerformance = `-- name: ListPlayerSeasonPerformance :many
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
    WHERE "match".season = $5::TEXT
    GROUP BY "player".id
    HAVING
        CASE
            WHEN $6::bool
            THEN EXISTS (
                SELECT 1
                FROM "club_player"
                WHERE
                    "club_player".player_id = "player".id AND
                    "club_player".club_id = $7::TEXT AND
                    "club_player".season = $5::TEXT
            )
            ELSE true
        END
), "player_ranked_total_stats" AS (
    SELECT
        player_total_stats.id, player_total_stats.total_goals, player_total_stats.total_assists, player_total_stats.total_yellow_cards, player_total_stats.total_red_cards, player_total_stats.total_own_goals, player_total_stats.appearances, player_total_stats.clean_sheets,
        $5::TEXT AS season,
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
SELECT id, total_goals, total_assists, total_yellow_cards, total_red_cards, total_own_goals, appearances, clean_sheets, season, goals_rank, assists_rank, clean_sheets_rank, fantasy_rank
FROM "player_ranked_total_stats"
WHERE
    CASE
        WHEN $1::bool
        THEN "player_ranked_total_stats".id = $2::INTEGER
        ELSE true
    END
ORDER BY
    CASE
        WHEN $3::TEXT = 'GOAL'
        THEN "player_ranked_total_stats".goals_rank
        WHEN $3::TEXT = 'ASSIST'
        THEN "player_ranked_total_stats".assists_rank
        WHEN $3::TEXT = 'CLEANSHEET'
        THEN "player_ranked_total_stats".clean_sheets_rank
        WHEN $3::TEXT = 'FANTASY'
        THEN "player_ranked_total_stats".fantasy_rank
        ELSE NULL
    END ASC
LIMIT $4::INTEGER
`

type ListPlayerSeasonPerformanceParams struct {
	FilterPlayerID bool
	PlayerID       int32
	OrderBy        string
	Limit          pgtype.Int4
	Season         string
	FilterClubID   bool
	ClubID         string
}

type ListPlayerSeasonPerformanceRow struct {
	ID               int32
	TotalGoals       int64
	TotalAssists     int64
	TotalYellowCards int64
	TotalRedCards    int64
	TotalOwnGoals    int64
	Appearances      int64
	CleanSheets      int32
	Season           string
	GoalsRank        int64
	AssistsRank      int64
	CleanSheetsRank  int64
	FantasyRank      int64
}

func (q *Queries) ListPlayerSeasonPerformance(ctx context.Context, arg ListPlayerSeasonPerformanceParams) ([]ListPlayerSeasonPerformanceRow, error) {
	rows, err := q.db.Query(ctx, listPlayerSeasonPerformance,
		arg.FilterPlayerID,
		arg.PlayerID,
		arg.OrderBy,
		arg.Limit,
		arg.Season,
		arg.FilterClubID,
		arg.ClubID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPlayerSeasonPerformanceRow
	for rows.Next() {
		var i ListPlayerSeasonPerformanceRow
		if err := rows.Scan(
			&i.ID,
			&i.TotalGoals,
			&i.TotalAssists,
			&i.TotalYellowCards,
			&i.TotalRedCards,
			&i.TotalOwnGoals,
			&i.Appearances,
			&i.CleanSheets,
			&i.Season,
			&i.GoalsRank,
			&i.AssistsRank,
			&i.CleanSheetsRank,
			&i.FantasyRank,
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

const updatePlayerByID = `-- name: UpdatePlayerByID :one
UPDATE "player" SET
    firstname = $1,
    lastname = $2,
    dob = $3,
    height = $4,
    nationality = $5,
    position = $6,
    image = $7
WHERE "player".id = $8
RETURNING id, firstname, lastname, dob, height, nationality, position, image
`

type UpdatePlayerByIDParams struct {
	Firstname   string
	Lastname    string
	Dob         pgtype.Timestamp
	Height      int16
	Nationality string
	Position    PlayerPosition
	Image       pgtype.Text
	ID          int32
}

func (q *Queries) UpdatePlayerByID(ctx context.Context, arg UpdatePlayerByIDParams) (Player, error) {
	row := q.db.QueryRow(ctx, updatePlayerByID,
		arg.Firstname,
		arg.Lastname,
		arg.Dob,
		arg.Height,
		arg.Nationality,
		arg.Position,
		arg.Image,
		arg.ID,
	)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Dob,
		&i.Height,
		&i.Nationality,
		&i.Position,
		&i.Image,
	)
	return i, err
}
