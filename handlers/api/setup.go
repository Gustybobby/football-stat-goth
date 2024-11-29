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
	r_api.Patch("/password", handlers.Make(HandleUpdatePassword, repo))

	r_api.Post("/fantasy/teams", handlers.Make(HandleCreateFantasyTeam, repo))
	r_api.Post("/fantasy/players/{playerID}", handlers.Make(HandleBuyFantasyPlayer, repo))
	r_api.Delete("/fantasy/players/{playerID}", handlers.Make(HandleSellFantasyPlayer, repo))

	r_api.Route("/users/{username}", func(r_api_user chi.Router) {
		SetupUserApiRoutes(r_api_user, repo)
	})

	r_api.Route("/admin", func(r_api_admin chi.Router) {
		SetupAdminApiRoutes(r_api_admin, repo)
	})
}

// prefix path '/api/users/{username}'
func SetupUserApiRoutes(r_api_user chi.Router, repo *repos.Repository) {
	r_api_user.Use(plmiddleware.AuthSessionUser)

	r_api_user.Patch("/", handlers.Make(HandleUpdateUser, repo))
}

// prefix path '/api/admin'
func SetupAdminApiRoutes(r_api_admin chi.Router, repo *repos.Repository) {
	r_api_admin.Use(plmiddleware.AuthAdmin)

	r_api_admin.Post("/players", handlers.Make(HandleCreatePlayer, repo))
	r_api_admin.Patch("/players/{playerID}", handlers.Make(HandleUpdatePlayer, repo))
	r_api_admin.Delete("/players/{playerID}", handlers.Make(HandleDeletePlayer, repo))

	r_api_admin.Post("/lineups/{lineupID}/lineup_players", handlers.Make(HandleCreateLineupPlayer, repo))
	r_api_admin.Patch("/lineups/{lineupID}/lineup_players/{playerID}", handlers.Make(HandleUpdateLineupPlayer, repo))
	r_api_admin.Delete("/lineups/{lineupID}/lineup_players/{playerID}", handlers.Make(HandleDeleteLineupPlayer, repo))
}
