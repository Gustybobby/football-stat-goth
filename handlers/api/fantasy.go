package api

import (
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/pltime"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func HandleCreateFantasyTeam(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	// player_0-10
	var added_player_id []int32
	var fantasy_players []queries.ListFantasyPlayersRow

	for i := range 11 {
		id, err := strconv.Atoi(r.FormValue("player_" + strconv.Itoa(i)))
		if err != nil {
			return err
		}

		added_player_id = append(added_player_id, int32(id))
	}

	fantasy_players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		MinCost:               1,
		AvgCost:               9,
		FilterFantasyPlayerID: true,
		FantasyPlayerIds:      added_player_id,
		Season:                pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return err
	}

	var fantasy_transaction []queries.InsertFantasyTransacionParams

	for i, fantasy_player := range fantasy_players {
		FantasyTeamID, err := strconv.Atoi(r.FormValue("player_" + strconv.Itoa(i) + "FantasyTeamID"))
		if err != nil {
			return err
		}

		FantasyPlayerID, err := strconv.Atoi(r.FormValue("player_" + strconv.Itoa(i) + "FantasyPlayerID"))
		if err != nil {
			return err
		}

		n := queries.InsertFantasyTransacionParams{
			Cost:            fantasy_player.Cost,
			Type:            queries.FantasyTransactionTypeBUY,
			FantasyTeamID:   pgtype.Int4{Int32: int32(FantasyTeamID), Valid: true},
			FantasyPlayerID: pgtype.Int4{Int32: int32(FantasyPlayerID), Valid: true},
		}

		fantasy_transaction = append(fantasy_transaction, n)
	}

	repo.Queries.InsertFantasyTransacion(repo.Ctx, fantasy_transaction)

	return nil
}
