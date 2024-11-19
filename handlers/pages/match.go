package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandleMatchPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	matchID, err := strconv.Atoi(chi.URLParam(r, "matchID"))
	if err != nil {
		return err
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

	match, err := repo.Queries.FindMatchByID(repo.Ctx, int32(matchID))
	if err != nil {
		return err
	}

	homeLineupPlayers, err := repo.Queries.ListLineupPlayersByLineupID(repo.Ctx, match.HomeLineupID)
	if err != nil {
		return err
	}

	awayLineupPlayers, err := repo.Queries.ListLineupPlayersByLineupID(repo.Ctx, match.AwayLineupID)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Match(fixtures, match, homeLineupPlayers, awayLineupPlayers))
}
