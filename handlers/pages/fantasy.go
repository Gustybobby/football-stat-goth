package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
)

func HandleFantasyPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	return handlers.Render(w, r, views.Fantasy())
}
