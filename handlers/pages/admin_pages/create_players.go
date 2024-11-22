package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views/admin/admin_views"
	"net/http"
)

func HandleAdminCreatePlayersPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user, err := plauth.Auth(w, r, repo)
	if err != nil {
		return err
	}
	if user.Role != queries.UserRoleADMIN {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	return handlers.Render(w, r, admin_views.CreatePlayers())
}
