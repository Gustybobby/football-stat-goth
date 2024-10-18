package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/models"
	"football-stat-goth/repos"
	"football-stat-goth/views"
	"net/http"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	var teams []models.Team
	repo.DB.Find(&teams)

	return handlers.Render(w, r, views.Home(teams))
}
