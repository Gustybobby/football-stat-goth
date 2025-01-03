package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
)

func HandleClubsPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	clubs, err := repo.Queries.ListClubsOrderByNameAsc(repo.Ctx)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Clubs(user, clubs))
}
