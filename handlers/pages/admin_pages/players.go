package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views/admin/admin_views"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleAdminCreatePlayersPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	return handlers.Render(w, r, admin_views.CreatePlayers())
}

func HandleAdminEditPlayersPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	playerID, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	player, err := repo.Queries.FindPlayerByID(repo.Ctx, int32(playerID))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return err
	}

	return handlers.Render(w, r, admin_views.EditPlayers(player))
}
