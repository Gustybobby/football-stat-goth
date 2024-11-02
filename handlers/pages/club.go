package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
)

func HandleClubPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	clubID := chi.URLParam(r, "clubID")

	db, conn, ctx, err := repo.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	club, err := db.FindClubByID(ctx, clubID)
	if err != nil {
		return err
	}

	fixtures, err := db.ListMatchesWithClubsAndGoals(ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: true,
		ClubID:       clubID,
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	matches, err := db.ListMatchesWithClubsAndGoals(ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: true,
		ClubID:       clubID,
		IsFinished:   true,
		Order:        "DESC",
	})
	if err != nil {
		return err
	}

	standings, err := db.ListClubStandings(ctx)
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(standings, func(club queries.ListClubStandingsRow) bool {
		return club.ID == clubID
	})

	averageStats, err := db.ClubAverageStatistics(ctx, clubID)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Club(club, fixtures, matches, standings[idx], idx+1, averageStats))
}
