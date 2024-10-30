package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HandleClubPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	clubID := chi.URLParam(r, "clubID")
	club, err := repos.FindClub(clubID, repo)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Club(*club))
}
