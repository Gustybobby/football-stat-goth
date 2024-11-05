package api

import (
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func HandleCreatePlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	db, conn, ctx, err := repo.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	no, err := strconv.Atoi(r.FormValue("no"))
	if err != nil {
		return err
	}

	dob, err := time.Parse("2006-01-02", r.FormValue("dob"))
	if err != nil {
		return err
	}

	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil {
		return err
	}

	db.CreatePlayer(ctx, queries.CreatePlayerParams{
		ClubID:      pgtype.Text{String: r.FormValue("club_id"), Valid: true},
		No:          int16(no),
		Firstname:   r.FormValue("firstname"),
		Lastname:    r.FormValue("lastname"),
		Dob:         pgtype.Timestamp{Time: dob, Valid: true},
		Height:      int16(height),
		Nationality: r.FormValue("nationality"),
		Position:    queries.PlayerPosition(r.FormValue("position")),
		Image:       pgtype.Text{String: r.FormValue("image"), Valid: true},
	})

	return nil
}
