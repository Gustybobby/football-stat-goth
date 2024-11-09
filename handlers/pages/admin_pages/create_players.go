package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views/admin/admin_views"
	"net/http"
)

func HandleAdminCreatePlayersPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	return handlers.Render(w, r, admin_views.CreatePlayers())
}