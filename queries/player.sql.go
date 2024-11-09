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
    club_id,
    no,
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
    $7,
    $8,
    $9
)
RETURNING id, club_id, no, firstname, lastname, dob, height, nationality, position, image
`

type CreatePlayerParams struct {
	ClubID      pgtype.Text
	No          int16
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
		arg.ClubID,
		arg.No,
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
		&i.ClubID,
		&i.No,
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