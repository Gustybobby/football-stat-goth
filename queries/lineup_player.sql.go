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
    no,
    position_no,
    position
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING lineup_id, player_id, no, position_no, position
`

type CreateLineupPlayerParams struct {
	LineupID   int32
	PlayerID   int32
	No         int16
	PositionNo int16
	Position   PlayerPosition
}

func (q *Queries) CreateLineupPlayer(ctx context.Context, arg CreateLineupPlayerParams) (LineupPlayer, error) {
	row := q.db.QueryRow(ctx, createLineupPlayer,
		arg.LineupID,
		arg.PlayerID,
		arg.No,
		arg.PositionNo,
		arg.Position,
	)
	var i LineupPlayer
	err := row.Scan(
		&i.LineupID,
		&i.PlayerID,
		&i.No,
		&i.PositionNo,
		&i.Position,
	)
	return i, err
}

const deleteLineupPlayer = `-- name: DeleteLineupPlayer :exec
DELETE FROM "lineup_player"
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".player_id = $2
`

type DeleteLineupPlayerParams struct {
	LineupID int32
	PlayerID int32
}

func (q *Queries) DeleteLineupPlayer(ctx context.Context, arg DeleteLineupPlayerParams) error {
	_, err := q.db.Exec(ctx, deleteLineupPlayer, arg.LineupID, arg.PlayerID)
	return err
}

const findLineupPlayerByLineupIDAndPositionNo = `-- name: FindLineupPlayerByLineupIDAndPositionNo :one
SELECT
    lineup_player.lineup_id, lineup_player.player_id, lineup_player.no, lineup_player.position_no, lineup_player.position,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "player"
ON "lineup_player".player_id = "player".id
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".position_no = $2
LIMIT 1
`

type FindLineupPlayerByLineupIDAndPositionNoParams struct {
	LineupID   int32
	PositionNo int16
}

type FindLineupPlayerByLineupIDAndPositionNoRow struct {
	LineupID   int32
	PlayerID   int32
	No         int16
	PositionNo int16
	Position   PlayerPosition
	Firstname  string
	Lastname   string
	Image      pgtype.Text
}

func (q *Queries) FindLineupPlayerByLineupIDAndPositionNo(ctx context.Context, arg FindLineupPlayerByLineupIDAndPositionNoParams) (FindLineupPlayerByLineupIDAndPositionNoRow, error) {
	row := q.db.QueryRow(ctx, findLineupPlayerByLineupIDAndPositionNo, arg.LineupID, arg.PositionNo)
	var i FindLineupPlayerByLineupIDAndPositionNoRow
	err := row.Scan(
		&i.LineupID,
		&i.PlayerID,
		&i.No,
		&i.PositionNo,
		&i.Position,
		&i.Firstname,
		&i.Lastname,
		&i.Image,
	)
	return i, err
}

const listLineupPlayersByLineupID = `-- name: ListLineupPlayersByLineupID :many
SELECT
    lineup_player.lineup_id, lineup_player.player_id, lineup_player.no, lineup_player.position_no, lineup_player.position,
    "player".firstname,
    "player".lastname,
    "player".image
FROM "lineup_player"
INNER JOIN "player"
ON "lineup_player".player_id = "player".id
WHERE "lineup_player".lineup_id = $1
ORDER BY "lineup_player".no ASC
`

type ListLineupPlayersByLineupIDRow struct {
	LineupID   int32
	PlayerID   int32
	No         int16
	PositionNo int16
	Position   PlayerPosition
	Firstname  string
	Lastname   string
	Image      pgtype.Text
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
			&i.No,
			&i.PositionNo,
			&i.Position,
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

const updateLineupPlayer = `-- name: UpdateLineupPlayer :one
UPDATE lineup_player SET
    position_no = COALESCE($3, position_no),
    position = COALESCE($4, position)
WHERE
    "lineup_player".lineup_id = $1 AND
    "lineup_player".player_id = $2
RETURNING lineup_id, player_id, no, position_no, position
`

type UpdateLineupPlayerParams struct {
	LineupID   int32
	PlayerID   int32
	PositionNo pgtype.Int2
	Position   NullPlayerPosition
}

func (q *Queries) UpdateLineupPlayer(ctx context.Context, arg UpdateLineupPlayerParams) (LineupPlayer, error) {
	row := q.db.QueryRow(ctx, updateLineupPlayer,
		arg.LineupID,
		arg.PlayerID,
		arg.PositionNo,
		arg.Position,
	)
	var i LineupPlayer
	err := row.Scan(
		&i.LineupID,
		&i.PlayerID,
		&i.No,
		&i.PositionNo,
		&i.Position,
	)
	return i, err
}
