package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
)

func HandleStandingsPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	db, conn, ctx, err := repo.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	clubs, err := db.ListClubStandings(ctx)
	if err != nil {
		return err
	}
	return handlers.Render(w, r, views.Standings(clubs))
}
