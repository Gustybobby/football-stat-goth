package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views/admin/admin_views"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleAdminEditMatchPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	matchID, err := strconv.Atoi(chi.URLParam(r, "matchID"))
	if err != nil {
		return err
	}

	match, err := repo.Queries.FindMatchByID(repo.Ctx, int32(matchID))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return err
	}

	events, err := repo.Queries.ListLineupEventsByMatchID(repo.Ctx, int32(matchID))
	if err != nil {
		return err
	}

	homeLineupPlayers, err := repo.Queries.ListLineupPlayersByLineupID(repo.Ctx, match.HomeLineupID)
	if err != nil {
		return err
	}

	awayLineupPlayers, err := repo.Queries.ListLineupPlayersByLineupID(repo.Ctx, match.AwayLineupID)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, admin_views.EditLineups(match, events, homeLineupPlayers, awayLineupPlayers))
}
