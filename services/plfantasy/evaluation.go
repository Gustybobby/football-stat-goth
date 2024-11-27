package plfantasy

import (
	"football-stat-goth/repos"
	"math"
)

type PlayerCost func(rank int) int

func GetPlayerCostFunc(season string, min_cost float32, budget float32, repo *repos.Repository) (PlayerCost, error) {
	rank_stats, err := repo.Queries.FindPlayersRankStats(repo.Ctx, season)
	if err != nil {
		return nil, err
	}

	avg_player_cost := budget / 11
	rate := (avg_player_cost - min_cost) / (float32(rank_stats.RankSum)/float32(rank_stats.PlayerCount) - 1)

	return func(rank int) int {
		return int(math.Round(float64(min_cost + rate*float32(rank_stats.MaxRank-int32(rank)))))
	}, nil
}
