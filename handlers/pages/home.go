package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user, err := plauth.Auth(w, r, repo)
	if err != nil {
		return err
	}

	fixtures, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		ClubID:       "",
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Home(user, fixtures))
}
