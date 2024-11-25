package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandlePlayerPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	player_id, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	player, err := repo.Queries.FindPlayerByID(repo.Ctx, int32(player_id))
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Player(user, player))
}
