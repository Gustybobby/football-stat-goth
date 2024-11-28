package api

import (
	"errors"
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/components/fantasy_components"
	"net/http"
	"strconv"
	"strings"
)

func HandleCreateFantasyTeam(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

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

	cost := 0
	for _, fantasy_player := range fantasy_players {
		cost += int(fantasy_player.Cost)
	}

	if r.Form.Get("submit_team") == "submit" {
		if cost > 100 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}

		fantasy_team, err := repo.Queries.CreateFantasyTeam(repo.Ctx, queries.CreateFantasyTeamParams{
			Username: user.Username,
			Season:   pltime.GetCurrentSeasonString(),
			Budget:   100,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return err
		}

		var fantasy_transactions []queries.CreateFantasyTransactionParams
		for _, fantasy_player := range fantasy_players {

			fantasy_transactions = append(fantasy_transactions, queries.CreateFantasyTransactionParams{
				Cost:            fantasy_player.Cost,
				Type:            queries.FantasyTransactionTypeBUY,
				FantasyTeamID:   fantasy_team.ID,
				FantasyPlayerID: fantasy_player.ID,
			})
		}

		repo.Queries.CreateFantasyTransaction(repo.Ctx, fantasy_transactions)
	}

	players_params, cost, err := GetFantasyTeamFieldParams(fantasy_players)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	return handlers.Render(w, r, fantasy_components.FantasyTeamField(*players_params, cost))
}

func GetFantasyTeamFieldParams(fantasy_players []queries.ListFantasyPlayersRow) (*fantasy_components.FantasyTeamFieldPlayersParams, int, error) {
	gk_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionGK, fantasy_players)
	if len(gk_fantasy_players) > 1 {
		return nil, 0, errors.New("GK players count exceed maximum")
	}

	def_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionDEF, fantasy_players)
	if len(def_fantasy_players) > 4 {
		return nil, 0, errors.New("DEF players count exceed maximum")
	}

	mfd_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionMFD, fantasy_players)
	if len(mfd_fantasy_players) > 4 {
		return nil, 0, errors.New("MFD players count exceed maximum")
	}

	fwd_fantasy_players := filterFantasyPlayersByPosition(queries.PlayerPositionFWD, fantasy_players)
	if len(fwd_fantasy_players) > 2 {
		return nil, 0, errors.New("FWD players count exceed maximum")
	}

	cost := 0
	for _, fantasy_player := range fantasy_players {
		cost += int(fantasy_player.Cost)
	}

	return &fantasy_components.FantasyTeamFieldPlayersParams{
		GK:  gk_fantasy_players,
		DEF: def_fantasy_players,
		MFD: mfd_fantasy_players,
		FWD: fwd_fantasy_players,
	}, cost, nil
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
