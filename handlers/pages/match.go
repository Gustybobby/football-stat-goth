package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleMatchPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	matchID, err := strconv.Atoi(chi.URLParam(r, "matchID"))
	if err != nil {
		return err
	}

	db, conn, ctx, err := repo.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	fixtures, err := db.ListMatchesWithClubsAndGoals(ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		ClubID:       "",
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	match, err := db.FindMatchByID(ctx, int32(matchID))
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Match(fixtures, match))
}
