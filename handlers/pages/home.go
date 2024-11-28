package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/plperformance"
	"football-stat-goth/views"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	fixtures, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		FilterWeek:   false,
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

	last_week, err := repo.Queries.FindLatestFinishedMatchweek(repo.Ctx)
	if err != nil {
		return err
	}

	last_week_matches, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		FilterWeek:   true,
		Week:         int32(last_week),
		IsFinished:   true,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	top_goal_cards, err := plperformance.ListTopPlayerCards("GOAL", "", 3, repo)
	if err != nil {
		return err
	}

	top_assist_cards, err := plperformance.ListTopPlayerCards("ASSIST", "", 3, repo)
	if err != nil {
		return err
	}

	top_cleansheet_cards, err := plperformance.ListTopPlayerCards("CLEANSHEET", "", 3, repo)
	if err != nil {
		return err
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
