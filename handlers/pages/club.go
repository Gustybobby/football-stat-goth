package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HandleClubPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	clubID := chi.URLParam(r, "clubID")

	db, conn, ctx, err := repo.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	club, err := db.FindClubByID(ctx, clubID)
	if err != nil {
		return err
	}

	fixtures, err := db.ListMatchesWithClubsAndGoals(ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: true,
		ClubID:       clubID,
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	matches, err := db.ListMatchesWithClubsAndGoals(ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: true,
		ClubID:       clubID,
		IsFinished:   true,
		Order:        "DESC",
	})
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Club(club, fixtures, matches))
}
