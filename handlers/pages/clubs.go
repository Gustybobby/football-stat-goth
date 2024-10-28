package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
)

func HandleClubsPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	clubs, err := repos.FindClubsByNameAsc(repo)
	if err != nil {
		return err
	}

	matches, err := repos.FindFixtureMatches(repo)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Clubs(matches, clubs))
}
