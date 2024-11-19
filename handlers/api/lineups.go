package api

import (
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleCreateLineupPlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	lineupID, err := strconv.Atoi(chi.URLParam(r, "lineupID"))
	if err != nil {
		return err
	}

	db, conn, ctx, err := repo.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	no, err := strconv.Atoi(r.FormValue("no"))
	if err != nil {
		return err
	}

	player_id, err := db.FindPlayerIDByClubAndNo(ctx, queries.FindPlayerIDByClubAndNoParams{
		ClubID: pgtype.Text{String: r.FormValue("club_id"), Valid: true},
		No:     int16(no),
	})
	if err != nil {
		return err
	}

	position_no, err := strconv.Atoi(r.FormValue("position_no"))
	if err != nil {
		return err
	}

	db.CreateLineupPlayer(ctx, queries.CreateLineupPlayerParams{
		LineupID:   int32(lineupID),
		PlayerID:   player_id,
		PositionNo: int16(position_no),
		Position:   queries.PlayerPosition(r.FormValue("position")),
	})

	return nil
}
