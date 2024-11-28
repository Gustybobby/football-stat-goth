package api

import (
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleCreatePlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	dob, err := time.Parse("2006-01-02", r.FormValue("dob"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	player, err := repo.Queries.CreatePlayer(repo.Ctx, queries.CreatePlayerParams{
		Firstname:   r.FormValue("firstname"),
		Lastname:    r.FormValue("lastname"),
		Dob:         pgtype.Timestamp{Time: dob, Valid: true},
		Height:      int16(height),
		Nationality: r.FormValue("nationality"),
		Position:    queries.PlayerPosition(r.FormValue("position")),
		Image:       pgtype.Text{String: r.FormValue("image"), Valid: true},
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	w.Header().Add("HX-Redirect", "/admin/players/"+strconv.Itoa(int(player.ID)))
	return nil
}

func HandleUpdatePlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	playerID, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	dob, err := time.Parse("2006-01-02", r.FormValue("dob"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	height, err := strconv.Atoi(r.FormValue("height"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	repo.Queries.UpdatePlayerByID(repo.Ctx, queries.UpdatePlayerByIDParams{
		ID:          int32(playerID),
		Firstname:   r.FormValue("firstname"),
		Lastname:    r.FormValue("lastname"),
		Dob:         pgtype.Timestamp{Time: dob, Valid: true},
		Height:      int16(height),
		Nationality: r.FormValue("nationality"),
		Position:    queries.PlayerPosition(r.FormValue("position")),
		Image:       pgtype.Text{String: r.FormValue("image"), Valid: true},
	})

	w.Header().Add("HX-Refresh", "true")
	return nil
}

func HandleDeletePlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	playerID, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	if r.FormValue("confirm") != strconv.Itoa(playerID) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	repo.Queries.DeletePlayerByID(repo.Ctx, int32(playerID))

	w.Header().Add("HX-Redirect", "/admin/players")
	return nil
}
