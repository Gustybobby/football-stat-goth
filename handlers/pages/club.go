package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/plperformance"
	"football-stat-goth/views"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
)

func HandleClubPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	club_id := chi.URLParam(r, "clubID")

	club, err := repo.Queries.FindClubByID(repo.Ctx, club_id)
	if err != nil {
		return err
	}

	fixtures, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: true,
		ClubID:       club_id,
		FilterWeek:   false,
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	matches, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: true,
		ClubID:       club_id,
		FilterWeek:   false,
		IsFinished:   true,
		Order:        "DESC",
	})
	if err != nil {
		return err
	}

	standings, err := repo.Queries.ListClubStandings(repo.Ctx)
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(standings, func(club queries.ListClubStandingsRow) bool {
		return club.ID == club_id
	})

	averageStats, err := repo.Queries.ClubAverageStatistics(repo.Ctx, club_id)
	if err != nil {
		return err
	}

	top_overall_cards, err := plperformance.ListTopPlayerCards("FANTASY", club_id, 3, repo)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Club(user, fixtures, club, matches, standings[idx], idx+1, averageStats, top_overall_cards))
}
