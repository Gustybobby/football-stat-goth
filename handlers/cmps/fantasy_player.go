package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/components/playercard_components"
	"net/http"
	"strconv"
)

func HandleFantasyPlayerCard(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	isBlank := r.URL.Query().Get("blank") == "true"

	if isBlank {
		return handlers.Render(w, r, playercard_components.BlankPlayerCard())
	}

	fantasy_player_id, err := strconv.Atoi(r.URL.Query().Get("fantasy_player_id"))
	if err != nil {
		return err
	}

	total_count, err := strconv.Atoi(r.FormValue(r.URL.Query().Get("position") + "_count"))
	if err != nil {
		return err
	}

	pos_max := map[string]int{
		"GK":  1,
		"DEF": 4,
		"MFD": 4,
		"FWD": 2,
	}
	if total_count > pos_max[r.URL.Query().Get("position")] {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	fantasy_players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		FilterFantasyPlayerID: true,
		FantasyPlayerID:       int32(fantasy_player_id),
		MinCost:               1,
		AvgCost:               9,
		Season:                pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return err
	}

	return handlers.Render(w, r, playercard_components.PlayerCardImg(fantasy_players[0]))
}
