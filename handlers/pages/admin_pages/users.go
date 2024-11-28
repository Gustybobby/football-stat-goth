package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views/admin/admin_views"
	"net/http"
)

func HandleAdminUsersPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	users, err := repo.Queries.ListUsers(repo.Ctx)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, admin_views.Users(users))
}
