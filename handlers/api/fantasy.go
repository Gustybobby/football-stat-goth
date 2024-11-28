package api

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/components/fantasy_components"
	"net/http"
	"strconv"
	"strings"
)

func HandleCreateFantasyTeam(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	r.ParseForm()

	var fantasy_player_ids []int32
	for key := range r.Form {
		if key == "submit_team" {
			continue
		}

		fantasy_player_id, err := strconv.Atoi(strings.Split(key, "_")[2])
		if err != nil {
			return err
		}

		fantasy_player_ids = append(fantasy_player_ids, int32(fantasy_player_id))
	}

	fantasy_players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		MinCost:               1,
		AvgCost:               9,
		FilterFantasyPlayerID: true,
		FantasyPlayerIds:      fantasy_player_ids,
		Season:                pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return err
	}

	if r.Form.Get("submit_team") == "submit" {
		var fantasy_transactions []queries.InsertFantasyTransacionParams
		for i, fantasy_player := range fantasy_players {
			FantasyTeamID, err := strconv.Atoi(r.FormValue("player_" + strconv.Itoa(i) + "FantasyTeamID"))
			if err != nil {
				return err
			}

			transaction := queries.InsertFantasyTransacionParams{
				Cost:            fantasy_player.Cost,
				Type:            queries.FantasyTransactionTypeBUY,
				FantasyTeamID:   int32(FantasyTeamID),
				FantasyPlayerID: fantasy_player.ID,
			}

			fantasy_transactions = append(fantasy_transactions, transaction)
		}

		repo.Queries.InsertFantasyTransacion(repo.Ctx, fantasy_transactions)
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
