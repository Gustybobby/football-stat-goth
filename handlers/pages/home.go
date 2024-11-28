package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/plperformance"
	"football-stat-goth/views"
	"football-stat-goth/views/components"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	ch_fixtures := make(chan []queries.ListMatchesWithClubsAndGoalsRow)
	ch_clubs := make(chan []queries.ListClubStandingsRow)
	ch_last_week := make(chan int16)
	ch_last_week_matches := make(chan []queries.ListMatchesWithClubsAndGoalsRow)
	ch_top_goal_cards := make(chan []components.PlayerPerformanceCardParams)
	ch_top_assist_cards := make(chan []components.PlayerPerformanceCardParams)
	ch_top_cleansheet_cards := make(chan []components.PlayerPerformanceCardParams)
	ch_err := make(chan error)

	var fixtures []queries.ListMatchesWithClubsAndGoalsRow
	var clubs []queries.ListClubStandingsRow
	var last_week int16
	var last_week_matches []queries.ListMatchesWithClubsAndGoalsRow
	var top_goal_cards []components.PlayerPerformanceCardParams
	var top_assist_cards []components.PlayerPerformanceCardParams
	var top_cleansheet_cards []components.PlayerPerformanceCardParams
	var err error

	go func() {
		rows, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
			FilterClubID: false,
			FilterWeek:   false,
			IsFinished:   false,
			Order:        "ASC",
		})
		ch_fixtures <- rows
		ch_err <- err
	}()

	go func() {
		rows, err := repo.Queries.ListClubStandings(repo.Ctx)
		ch_clubs <- rows
		ch_err <- err
	}()

	go func() {
		row, err := repo.Queries.FindLatestFinishedMatchweek(repo.Ctx)
		ch_last_week <- row
		ch_err <- err

		rows, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
			FilterClubID: false,
			FilterWeek:   true,
			Week:         int32(row),
			IsFinished:   true,
			Order:        "ASC",
		})
		ch_last_week_matches <- rows
		ch_err <- err
	}()

	go func() {
		rows, err := plperformance.ListTopPlayerCards("GOAL", "", 3, repo)
		ch_top_goal_cards <- rows
		ch_err <- err
	}()

	go func() {
		rows, err := plperformance.ListTopPlayerCards("ASSIST", "", 3, repo)
		ch_top_assist_cards <- rows
		ch_err <- err
	}()

	go func() {
		rows, err := plperformance.ListTopPlayerCards("CLEANSHEET", "", 3, repo)
		ch_top_cleansheet_cards <- rows
		ch_err <- err
	}()

	for i := 0; i < 14; i++ {
		select {
		case fixtures = <-ch_fixtures:
		case clubs = <-ch_clubs:
		case last_week = <-ch_last_week:
		case last_week_matches = <-ch_last_week_matches:
		case top_goal_cards = <-ch_top_goal_cards:
		case top_assist_cards = <-ch_top_assist_cards:
		case top_cleansheet_cards = <-ch_top_cleansheet_cards:
		case err = <-ch_err:
		}
		if err != nil {
			return err
		}
	}

	return handlers.Render(w, r, views.Home(user, fixtures, clubs, views.MatchTableParams{
		Week:    last_week,
		Matches: last_week_matches,
	}, views.TopPlayersCardParams{
		Goal:       top_goal_cards,
		Assist:     top_assist_cards,
		CleanSheet: top_cleansheet_cards,
	}))
}
