package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views/admin/admin_views"
	"net/http"
)

func HandleAdminHomePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	return handlers.Render(w, r, admin_views.AdminHome(user))
}
