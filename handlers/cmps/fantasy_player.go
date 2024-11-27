package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/api"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/components/fantasy_components"
	"net/http"
	"strconv"
	"strings"
)

func HandleFantasyPlayersField(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	r.ParseForm()

	if r.Form.Get("submit_team") == "submit" {
		return api.HandleCreateFantasyTeam(w, r, repo)
	}

	var fantasy_player_ids []int32
	for key := range r.Form {
		fantasy_player_id, err := strconv.Atoi(strings.Split(key, "_")[2])
		if err != nil {
			return err
		}

		fantasy_player_ids = append(fantasy_player_ids, int32(fantasy_player_id))
	}

	fantasy_players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		FilterFantasyPlayerID: true,
		FantasyPlayerIds:      fantasy_player_ids,
		MinCost:               1,
		AvgCost:               9,
		Season:                pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return err
	}

	cost := 0
	for _, fantasy_player := range fantasy_players {
		cost += int(fantasy_player.Cost)
	}

	gk_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionGK, fantasy_players)
	if len(gk_fantasy_players) > 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	def_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionDEF, fantasy_players)
	if len(def_fantasy_players) > 4 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	mfd_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionMFD, fantasy_players)
	if len(mfd_fantasy_players) > 4 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	fwd_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionFWD, fantasy_players)
	if len(fwd_fantasy_players) > 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	return handlers.Render(w, r, fantasy_components.FantasyTeamField(fantasy_components.FantasyTeamFieldPlayersParams{
		GK:  gk_fantasy_players,
		DEF: def_fantasy_players,
		MFD: mfd_fantasy_players,
		FWD: fwd_fantasy_players,
	}, cost))
}

func filterFantasyPlayersByPosition(position queries.PlayerPosition, fantasy_players []queries.ListFantasyPlayersRow) []queries.ListFantasyPlayersRow {
	var filtered_fantasy_players []queries.ListFantasyPlayersRow
	for _, fantasy_player := range fantasy_players {
		if fantasy_player.Position == position {
			filtered_fantasy_players = append(filtered_fantasy_players, fantasy_player)
		}
	}
	return filtered_fantasy_players
}
