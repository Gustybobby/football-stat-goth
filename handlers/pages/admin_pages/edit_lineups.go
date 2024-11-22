package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views/admin/admin_views"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleAdminEditLineupsPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	matchID, err := strconv.Atoi(chi.URLParam(r, "matchID"))
	if err != nil {
		return err
	}

	user, err := plauth.Auth(w, r, repo)
	if err != nil {
		return err
	}
	if user.Role != queries.UserRoleADMIN {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	match, err := repo.Queries.FindMatchByID(repo.Ctx, int32(matchID))
	if err != nil {
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
