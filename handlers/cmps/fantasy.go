package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/api"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/plconstant"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/components/fantasy_components"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleFantasyPlayerDetailsComponent(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
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
		return err
	}
	fantasy_player := fantasy_players[0]

	fantasy_team, err := repo.Queries.FindFantasyTeamByUsernameSeason(repo.Ctx, queries.FindFantasyTeamByUsernameSeasonParams{
		Username: user.Username,
		Season:   pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return handlers.Render(w, r, fantasy_components.PlayerCardDetails(fantasy_player, false, 0, false))
	}

	transaction, err := repo.Queries.FindLastestTransaction(repo.Ctx, queries.FindLastestTransactionParams{
		FantasyTeamID:   fantasy_team.ID,
		FantasyPlayerID: fantasy_player.ID,
	})
	if err != nil || transaction.Type == queries.FantasyTransactionTypeSELL {
		is_valid_count, err := api.IsValidPositionCount(fantasy_player.Position, fantasy_team.ID, repo)
		if err != nil {
			return err
		}
		return handlers.Render(w, r, fantasy_components.PlayerCardDetails(fantasy_player, false, 0, is_valid_count && fantasy_team.Budget >= fantasy_player.Cost))
	}

	return handlers.Render(w, r, fantasy_components.PlayerCardDetails(fantasy_player, true, transaction.Cost, false))
}
