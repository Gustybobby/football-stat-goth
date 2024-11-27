package api

import (
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/pltime"
	"net/http"
	"strconv"
	"strings"
)

func HandleCreateFantasyTeam(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	r.ParseForm()

	var fantasy_player_ids []int32
	for key := range r.Form {
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

	return nil
}
