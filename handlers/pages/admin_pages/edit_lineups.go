package admin_pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
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

	db, conn, ctx, err := repo.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	match, err := db.FindMatchByID(ctx, int32(matchID))
	if err != nil {
		return err
	}

	homeLineupPlayers, err := db.ListLineupPlayersByLineupID(ctx, match.HomeLineupID)
	if err != nil {
		return err
	}

	awayLineupPlayers, err := db.ListLineupPlayersByLineupID(ctx, match.AwayLineupID)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, admin_views.EditLineups(match, homeLineupPlayers, awayLineupPlayers))
}
