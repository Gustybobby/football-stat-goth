package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"

	"github.com/go-chi/chi/v5"
)

func SetupComponentRoutes(r_cmps chi.Router, repo *repos.Repository) {
	r_cmps.Get("/lineup-players/form", handlers.Make(HandleLineupPlayerForm, repo))
}
