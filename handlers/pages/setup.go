package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/pages/admin_pages"
	"football-stat-goth/handlers/plmiddleware"
	"football-stat-goth/repos"

	"github.com/go-chi/chi/v5"
)

// prefix path '/'
func SetupPageRoutes(router *chi.Mux, repo *repos.Repository) {
	router.Get("/", handlers.Make(HandleHomePage, repo))
	router.Get("/standings", handlers.Make(HandleStandingsPage, repo))
	router.Get("/clubs", handlers.Make(HandleClubsPage, repo))
	router.Get("/fantasy", handlers.Make(HandleFantasyPage, repo))
	router.Get("/clubs/{clubID}", handlers.Make(HandleClubPage, repo))
	router.Get("/matches/{matchID}", handlers.Make(HandleMatchPage, repo))
	router.Get("/players", handlers.Make(HandlePlayersPage, repo))
	router.Get("/players/{playerID}", handlers.Make(HandlePlayerPage, repo))
	router.Get("/signup", handlers.Make(HandleSignupPage, repo))
	router.Get("/signin", handlers.Make(HandleSigninPage, repo))
	router.Get("/profile", handlers.Make(HandleProfilePage, repo))

	router.Route("/admin", func(r_admin chi.Router) {
		SetupAdminPageRoutes(r_admin, repo)
	})
}

// prefix path '/admin'
func SetupAdminPageRoutes(r_admin chi.Router, repo *repos.Repository) {
	r_admin.Use(plmiddleware.AuthAdmin)

	r_admin.Get("/", handlers.Make(admin_pages.HandleAdminHomePage, repo))

	r_admin.Get("/users", handlers.Make(admin_pages.HandleAdminUsersPage, repo))

	r_admin.Get("/players", handlers.Make(admin_pages.HandleAdminCreatePlayersPage, repo))
	r_admin.Get("/players/{playerID}", handlers.Make(admin_pages.HandleAdminEditPlayersPage, repo))

	r_admin.Get("/matches/{matchID}", handlers.Make(admin_pages.HandleAdminEditMatchPage, repo))
}
