package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
)

func HandlePlayersPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	return handlers.Render(w, r, views.Players(user))
}
