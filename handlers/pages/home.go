package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	matches, err := repos.FindFixtureMatches(repo)
	if err != nil {
		return err
	}
	return handlers.Render(w, r, views.Home(matches))
}
