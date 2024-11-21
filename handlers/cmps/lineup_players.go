package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/views/admin/admin_components/admin_lineup_components"
	"net/http"
	"strconv"
)

func HandleLineupPlayerForm(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	lineup_id, err := strconv.Atoi(r.URL.Query().Get("lineup_id"))
	if err != nil {
		return err
	}

	if r.URL.Query().Get("type") == "add" {
		position_no, err := strconv.Atoi(r.URL.Query().Get("position_no"))
		if err != nil {
			return err
		}
		if position_no >= 100 {
			return handlers.Render(w, r, admin_lineup_components.AddLineupSubstituteForm(admin_lineup_components.AddPlayerFormParams{
				LineupID:   lineup_id,
				PositionNo: r.URL.Query().Get("position_no"),
				ClubID:     r.URL.Query().Get("club_id"),
				Mirror:     r.URL.Query().Get("mirror"),
			}))
		}
		return handlers.Render(w, r, admin_lineup_components.AddLineupPlayerForm(admin_lineup_components.AddPlayerFormParams{
			LineupID:   lineup_id,
			PositionNo: r.URL.Query().Get("position_no"),
			ClubID:     r.URL.Query().Get("club_id"),
			Mirror:     r.URL.Query().Get("mirror"),
		}))
	}

	if r.URL.Query().Get("type") == "edit" {
		position_no, err := strconv.Atoi(r.URL.Query().Get("position_no"))
		if err != nil {
			return err
		}

		lineupPlayer, err := repo.Queries.FindLineupPlayerByLineupIDAndPositionNo(repo.Ctx, queries.FindLineupPlayerByLineupIDAndPositionNoParams{
			LineupID:   int32(lineup_id),
			PositionNo: int16(position_no),
		})
		if err != nil {
			return err
		}

		return handlers.Render(w, r, admin_lineup_components.EditLineupPlayerForm(lineupPlayer))
	}

	return nil
}
