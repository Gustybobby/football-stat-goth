package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views/admin/admin_components/admin_lineup_components"
	"net/http"
	"strconv"
)

func HandleAddLineupPlayerForm(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	lineup_id, err := strconv.Atoi(r.URL.Query().Get("lineup_id"))
	if err != nil {
		return err
	}

	return handlers.Render(w, r, admin_lineup_components.AddPlayerForm(admin_lineup_components.AddPlayerFormParams{
		LineupID:   lineup_id,
		PositionNo: r.URL.Query().Get("position_no"),
		ClubID:     r.URL.Query().Get("club_id"),
		Mirror:     r.URL.Query().Get("mirror"),
	}))
}
