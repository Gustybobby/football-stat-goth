package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/pages/admin_pages"
	"football-stat-goth/repos"

	"github.com/go-chi/chi/v5"
)

func SetupPageRoutes(router *chi.Mux, repo *repos.Repository) {
	router.Get("/", handlers.Make(HandleHomePage, repo))
	router.Get("/standings", handlers.Make(HandleStandingsPage, repo))
	router.Get("/clubs", handlers.Make(HandleClubsPage, repo))
	router.Get("/fantasy", handlers.Make(HandleFantasyPage, repo))
	router.Get("/clubs/{clubID}", handlers.Make(HandleClubPage, repo))
	router.Get("/matches/{matchID}", handlers.Make(HandleMatchPage, repo))
	router.Get("/signup", handlers.Make(HandleSignupPage, repo))
	router.Get("/signin", handlers.Make(HandleSigninPage, repo))

	router.Route("/admin", func(r_admin chi.Router) {
		SetupAdminPageRoutes(r_admin, repo)
	})
}

func SetupAdminPageRoutes(r_admin chi.Router, repo *repos.Repository) {
	r_admin.Get("/players/create", handlers.Make(admin_pages.HandleAdminCreatePlayersPage, repo))

	r_admin.Get("/matches/{matchID}/lineups", handlers.Make(admin_pages.HandleAdminEditLineupsPage, repo))
}
