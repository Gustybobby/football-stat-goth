package api

import (
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/plmiddleware"
	"football-stat-goth/repos"

	"github.com/go-chi/chi/v5"
)

// prefix path '/api'
func SetupApiRoutes(r_api chi.Router, repo *repos.Repository) {
	r_api.Post("/signup", handlers.Make(HandleSignup, repo))
	r_api.Post("/signin", handlers.Make(HandleSignin, repo))
	r_api.Delete("/signout", handlers.Make(HandleSignout, repo))
	r_api.Post("/fantasy/teams", handlers.Make(HandleCreateFantasyTeam, repo))
	r_api.Patch("/password", handlers.Make(HandleUpdatePassword, repo))

	r_api.Route("/admin", func(r_api_admin chi.Router) {
		SetupAdminApiRoutes(r_api_admin, repo)
	})
}

// prefix path '/api/admin'
func SetupAdminApiRoutes(r_api_admin chi.Router, repo *repos.Repository) {
	r_api_admin.Use(plmiddleware.AuthAdmin)

	r_api_admin.Post("/players", handlers.Make(HandleCreatePlayer, repo))

	r_api_admin.Post("/lineups/{lineupID}/lineup_players", handlers.Make(HandleCreateLineupPlayer, repo))
	r_api_admin.Patch("/lineups/{lineupID}/lineup_players/{playerID}", handlers.Make(HandleUpdateLineupPlayer, repo))
}
