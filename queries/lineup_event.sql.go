// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: lineup_event.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listLineupEventsByMatchID = `-- name: ListLineupEventsByMatchID :many
SELECT
    lineup_event.id, lineup_event.lineup_id, lineup_event.player_id1, lineup_event.player_id2, lineup_event.event, lineup_event.minutes, lineup_event.extra, lineup_event.after_half,
    "club_player1".no AS p1_no,
    "player1".firstname AS p1_firstname,
    "player1".lastname AS p1_lastname,
    "club_player2".no AS p2_no,
    "player2".firstname AS p2_firstname,
    "player2".lastname AS p2_lastname
FROM "lineup_event"
INNER JOIN "lineup"
ON "lineup_event".lineup_id = "lineup".id
INNER JOIN "match"
ON "lineup".id IN ("match".home_lineup_id, "match".away_lineup_id)
LEFT JOIN "club_player" AS "club_player1"
ON
    "lineup".club_id = "club_player1".club_id AND
    "lineup_event".player_id1 = "club_player1".player_id AND
    "match".season = "club_player1".season
LEFT JOIN "player" AS "player1"
ON "club_player1".player_id = "player1".id
LEFT JOIN "club_player" AS "club_player2"
ON
    "lineup".club_id = "club_player2".club_id AND
    "lineup_event".player_id2 = "club_player2".player_id AND
    "match".season = "club_player2".season
LEFT JOIN "player" AS "player2"
ON "club_player2".player_id = "player2".id
WHERE "match".id = $1
ORDER BY ("lineup_event".minutes + COALESCE("lineup_event".extra,0)) ASC
`

type ListLineupEventsByMatchIDRow struct {
	ID          int32
	LineupID    int32
	PlayerId1   pgtype.Int4
	PlayerId2   pgtype.Int4
	Event       EventType
	Minutes     int16
	Extra       pgtype.Int2
	AfterHalf   bool
	P1No        pgtype.Int2
	P1Firstname pgtype.Text
	P1Lastname  pgtype.Text
	P2No        pgtype.Int2
	P2Firstname pgtype.Text
	P2Lastname  pgtype.Text
}

func (q *Queries) ListLineupEventsByMatchID(ctx context.Context, id int32) ([]ListLineupEventsByMatchIDRow, error) {
	rows, err := q.db.Query(ctx, listLineupEventsByMatchID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListLineupEventsByMatchIDRow
	for rows.Next() {
		var i ListLineupEventsByMatchIDRow
		if err := rows.Scan(
			&i.ID,
			&i.LineupID,
			&i.PlayerId1,
			&i.PlayerId2,
			&i.Event,
			&i.Minutes,
			&i.Extra,
			&i.AfterHalf,
			&i.P1No,
			&i.P1Firstname,
			&i.P1Lastname,
			&i.P2No,
			&i.P2Firstname,
			&i.P2Lastname,
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
