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
	fantasy_player_id, err := strconv.Atoi(r.URL.Query().Get("fantasy_player_id"))
	if err != nil {
		return err
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
