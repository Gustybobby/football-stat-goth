// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: lineup_player.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createLineupPlayer = `-- name: CreateLineupPlayer :one
INSERT INTO "lineup_player" (
    lineup_id,
    player_id,
    position_no,
    position
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING lineup_id, player_id, position_no, position, goals, yellow_cards, red_cards
`

type CreateLineupPlayerParams struct {
	LineupID   int32
	PlayerID   int32
	PositionNo int16
	Position   PlayerPosition
}

func (q *Queries) CreateLineupPlayer(ctx context.Context, arg CreateLineupPlayerParams) (LineupPlayer, error) {
	row := q.db.QueryRow(ctx, createLineupPlayer,
		arg.LineupID,
		arg.PlayerID,
		arg.PositionNo,
		arg.Position,
	)
	var i LineupPlayer
	err := row.Scan(
		&i.LineupID,
		&i.PlayerID,
		&i.PositionNo,
		&i.Position,
		&i.Goals,
		&i.YellowCards,
		&i.RedCards,
	)
	return i, err
}

const listLineupPlayersByLineupID = `-- name: ListLineupPlayersByLineupID :many
SELECT
    lineup_player.lineup_id, lineup_player.player_id, lineup_player.position_no, lineup_player.position, lineup_player.goals, lineup_player.yellow_cards, lineup_player.red_cards,
    "player".no,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "player"
ON "lineup_player".player_id = "player".id
WHERE "lineup_player".lineup_id = $1
`

type ListLineupPlayersByLineupIDRow struct {
	LineupID    int32
	PlayerID    int32
	PositionNo  int16
	Position    PlayerPosition
	Goals       int16
	YellowCards int16
	RedCards    int16
	No          int16
	Firstname   string
	Lastname    string
	Image       pgtype.Text
}

func (q *Queries) ListLineupPlayersByLineupID(ctx context.Context, lineupID int32) ([]ListLineupPlayersByLineupIDRow, error) {
	rows, err := q.db.Query(ctx, listLineupPlayersByLineupID, lineupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLineupPlayersByLineupIDRow
	for rows.Next() {
		var i ListLineupPlayersByLineupIDRow
		if err := rows.Scan(
			&i.LineupID,
			&i.PlayerID,
			&i.PositionNo,
			&i.Position,
			&i.Goals,
			&i.YellowCards,
			&i.RedCards,
			&i.No,
			&i.Firstname,
			&i.Lastname,
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
