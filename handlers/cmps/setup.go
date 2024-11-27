package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/plmiddleware"
	"football-stat-goth/repos"

	"github.com/go-chi/chi/v5"
)

// prefix path '/cmps'
func SetupComponentRoutes(r_cmps chi.Router, repo *repos.Repository) {
	r_cmps.Get("/players-table", handlers.Make(HandlePlayersTable, repo))
	r_cmps.Post("/fantasy/player-card", handlers.Make(HandleFantasyPlayersField, repo))
	r_cmps.Route("/admin", func(r_cmps_admin chi.Router) {
		SetupAdminComponentRoutes(r_cmps_admin, repo)
	})
}

// prefix path '/cmps/admin'
func SetupAdminComponentRoutes(r_cmps_admin chi.Router, repo *repos.Repository) {
	r_cmps_admin.Use(plmiddleware.AuthAdmin)

	r_cmps_admin.Get("/lineup-players/form", handlers.Make(HandleLineupPlayerForm, repo))
}
