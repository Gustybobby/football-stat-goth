package plperformance

import (
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/components"

	"github.com/jackc/pgx/v5/pgtype"
)

func ListTopPlayerCards(order_by string, club_id string, limit int, repo *repos.Repository) ([]components.PlayerPerformanceCardParams, error) {
	top_players, err := repo.Queries.ListPlayerSeasonPerformance(repo.Ctx, queries.ListPlayerSeasonPerformanceParams{
		Season:         pltime.GetCurrentSeasonString(),
		FilterPlayerID: false,
		FilterClubID:   club_id != "",
		ClubID:         club_id,
		Limit:          pgtype.Int4{Int32: int32(limit), Valid: true},
		OrderBy:        order_by,
	})
	if err != nil {
		return nil, err
	}

	top_player_cards, err := listPerformanceCardParams(top_players, repo)
	if err != nil {
		return nil, err
	}

	return top_player_cards, nil
}

func listPerformanceCardParams(performances []queries.ListPlayerSeasonPerformanceRow, repo *repos.Repository) ([]components.PlayerPerformanceCardParams, error) {
	var cards []components.PlayerPerformanceCardParams
	for _, player_perf := range performances {
		player, err := repo.Queries.FindPlayerByID(repo.Ctx, player_perf.ID)
		if err != nil {
			return nil, err
		}

		club_players, err := repo.Queries.ListClubPlayerByPlayerID(repo.Ctx, player_perf.ID)
		if err != nil {
			return nil, err
		}

		cards = append(cards, components.PlayerPerformanceCardParams{
			Performance: player_perf,
			Player:      player,
			ClubPlayer:  club_players[0],
		})
	}
	return cards, nil
}
