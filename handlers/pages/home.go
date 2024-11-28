package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views"
	"football-stat-goth/views/components"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	fixtures, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		ClubID:       "",
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	clubs, err := repo.Queries.ListClubStandings(repo.Ctx)
	if err != nil {
		return err
	}

	top_goals, err := repo.Queries.ListPlayerSeasonPerformance(repo.Ctx, queries.ListPlayerSeasonPerformanceParams{
		Season:         pltime.GetCurrentSeasonString(),
		FilterPlayerID: false,
		FilterClubID:   false,
		Limit:          pgtype.Int4{Int32: 3, Valid: true},
		OrderBy:        "GOAL",
	})
	if err != nil {
		return err
	}

	top_goal_cards, err := listPerformanceCardParams(top_goals, repo)
	if err != nil {
		return err
	}

	top_assists, err := repo.Queries.ListPlayerSeasonPerformance(repo.Ctx, queries.ListPlayerSeasonPerformanceParams{
		Season:         pltime.GetCurrentSeasonString(),
		FilterPlayerID: false,
		FilterClubID:   false,
		Limit:          pgtype.Int4{Int32: 3, Valid: true},
		OrderBy:        "ASSIST",
	})
	if err != nil {
		return err
	}

	top_assist_cards, err := listPerformanceCardParams(top_assists, repo)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Home(user, fixtures, clubs, views.TopPlayersCardParams{Goal: top_goal_cards, Assist: top_assist_cards}))
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
