package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
)

func HandleClubsPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	clubs, err := repos.FindClubsWithNameAsc(repo)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Clubs(clubs))
}
