package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views"
	"net/http"
	"net/url"
)

func HandleFantasyPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Redirect(w, r, "/signin?redirectUrl="+url.QueryEscape("/fantasy"), http.StatusFound)
		return nil
	}

	fixtures, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		ClubID:       "",
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		FilterFantasyPlayerID: false,
		MinCost:               1,
		AvgCost:               9,
		Season:                pltime.GetCurrentSeasonString(),
	})

	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Fantasy(user, fixtures, players))
}
