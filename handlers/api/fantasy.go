package api

import (
	"errors"
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/plconstant"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/components/fantasy_components"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func HandleCreateFantasyTeam(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	_, err := repo.Queries.FindFantasyTeamByUsernameSeason(repo.Ctx, queries.FindFantasyTeamByUsernameSeasonParams{
		Username: user.Username,
		Season:   pltime.GetCurrentSeasonString(),
	})
	if err == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		MinCost:               plconstant.FantasyPlayerMinCost,
		AvgCost:               plconstant.FantasyPlayerAverageCost,
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
		if cost > plconstant.FantasyTeamMaxBudget || len(fantasy_players) != 11 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}

		fantasy_team, err := repo.Queries.CreateFantasyTeam(repo.Ctx, queries.CreateFantasyTeamParams{
			Username: user.Username,
			Season:   pltime.GetCurrentSeasonString(),
			Budget:   plconstant.FantasyTeamMaxBudget,
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
	players_params.HasTeam = false

	return handlers.Render(w, r, fantasy_components.FantasyTeamField(*players_params, plconstant.FantasyTeamMaxBudget-cost))
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
		GK:      gk_fantasy_players,
		DEF:     def_fantasy_players,
		MFD:     mfd_fantasy_players,
		FWD:     fwd_fantasy_players,
		HasTeam: len(fantasy_players) > 0,
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

func HandleSellFantasyPlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	fantasy_team, err := repo.Queries.FindFantasyTeamByUsernameSeason(repo.Ctx, queries.FindFantasyTeamByUsernameSeasonParams{
		Username: user.Username,
		Season:   pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return err
	}

	transaction, err := repo.Queries.FindLastestTransaction(repo.Ctx, queries.FindLastestTransactionParams{
		FantasyTeamID:   fantasy_team.ID,
		FantasyPlayerID: int32(playerID),
	})
	if err != nil || transaction.Type == queries.FantasyTransactionTypeSELL {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	fantasy_players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		MinCost:               plconstant.FantasyPlayerMinCost,
		AvgCost:               plconstant.FantasyPlayerAverageCost,
		FilterFantasyPlayerID: true,
		FantasyPlayerIds:      []int32{int32(playerID)},
		Season:                pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return err
	}
	fantasy_player := fantasy_players[0]

	repo.Queries.CreateFantasyTransaction(repo.Ctx, []queries.CreateFantasyTransactionParams{{
		Cost:            fantasy_player.Cost,
		Type:            queries.FantasyTransactionTypeSELL,
		FantasyTeamID:   fantasy_team.ID,
		FantasyPlayerID: fantasy_player.ID,
	}})

	w.Header().Add("HX-Refresh", "true")
	return nil
}

func HandleBuyFantasyPlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	fantasy_team, err := repo.Queries.FindFantasyTeamByUsernameSeason(repo.Ctx, queries.FindFantasyTeamByUsernameSeasonParams{
		Username: user.Username,
		Season:   pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return err
	}

	transaction, err := repo.Queries.FindLastestTransaction(repo.Ctx, queries.FindLastestTransactionParams{
		FantasyTeamID:   fantasy_team.ID,
		FantasyPlayerID: int32(playerID),
	})
	if err == nil && transaction.Type == queries.FantasyTransactionTypeBUY {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return err
	}

	fantasy_players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		MinCost:               plconstant.FantasyPlayerMinCost,
		AvgCost:               plconstant.FantasyPlayerAverageCost,
		FilterFantasyPlayerID: true,
		FantasyPlayerIds:      []int32{int32(playerID)},
		Season:                pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return err
	}
	fantasy_player := fantasy_players[0]

	is_valid_count, err := IsValidPositionCount(fantasy_player.Position, fantasy_team.ID, repo)
	if err != nil {
		return err
	}
	if !is_valid_count {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	repo.Queries.CreateFantasyTransaction(repo.Ctx, []queries.CreateFantasyTransactionParams{{
		Cost:            fantasy_player.Cost,
		Type:            queries.FantasyTransactionTypeBUY,
		FantasyTeamID:   fantasy_team.ID,
		FantasyPlayerID: fantasy_player.ID,
	}})

	w.Header().Add("HX-Refresh", "true")
	return nil
}

func IsValidPositionCount(position_to_add queries.PlayerPosition, fantasy_team_id int32, repo *repos.Repository) (bool, error) {
	players_count, err := repo.Queries.CountFantasyTeamPlayersByFantasyTeamID(repo.Ctx, fantasy_team_id)
	if err != nil {
		return false, err
	}

	if position_to_add == queries.PlayerPositionGK && players_count.GkCount < 1 {
		return true, nil
	}

	if position_to_add == queries.PlayerPositionDEF && players_count.DefCount < 4 {
		return true, nil
	}

	if position_to_add == queries.PlayerPositionMFD && players_count.MfdCount < 4 {
		return true, nil
	}

	if position_to_add == queries.PlayerPositionFWD && players_count.FwdCount < 2 {
		return true, nil
	}

	return false, nil
}
