package main

import (
	"embed"
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/api"
	"football-stat-goth/handlers/cmps"
	"football-stat-goth/handlers/pages"
	"football-stat-goth/handlers/pages/admin_pages"
	"football-stat-goth/repos"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

//go:embed public
var publicFS embed.FS

func SetupRoutes(router *chi.Mux, repo *repos.Repository) {
	router.Handle("/*", http.FileServerFS(publicFS))

	router.Get("/", handlers.Make(pages.HandleHomePage, repo))
	router.Get("/standings", handlers.Make(pages.HandleStandingsPage, repo))
	router.Get("/clubs", handlers.Make(pages.HandleClubsPage, repo))
	router.Get("/fantasy", handlers.Make(pages.HandleFantasyPage, repo))
	router.Get("/clubs/{clubID}", handlers.Make(pages.HandleClubPage, repo))
	router.Get("/matches/{matchID}", handlers.Make(pages.HandleMatchPage, repo))
	router.Get("/signup", handlers.Make(pages.HandleSignupPage, repo))

	router.Route("/admin", func(r chi.Router) {
		r.Get("/players/create", handlers.Make(admin_pages.HandleAdminCreatePlayersPage, repo))

		r.Get("/matches/{matchID}/lineups", handlers.Make(admin_pages.HandleAdminEditLineupsPage, repo))
	})

	router.Route("/api", func(r chi.Router) {
		r.Post("/signup", handlers.Make(api.HandleSignup, repo))
		r.Post("/players", handlers.Make(api.HandleCreatePlayer, repo))

		r.Post("/lineups/{lineupID}/lineup_players", handlers.Make(api.HandleCreateLineupPlayer, repo))
	})

	router.Route("/cmps", func(r chi.Router) {
		r.Get("/lineup-players/add-form", handlers.Make(cmps.HandleAddLineupPlayerForm, repo))
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(middleware.Logger)

	config := &repos.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	repo, err := repos.DbConnect(repos.Dsn(config))
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Conn.Close(repo.Ctx)

	SetupRoutes(router, repo)

	serverAddr := os.Getenv("SERVER_ADDR")

	slog.Info(`
	 ___  _                          _             
	| . \| |   ___  _ _  _ _ _  ___ | |__ ___  _ _ 
	|  _/| |_ [_] || | || ' ' |[_] || / // ._]| '_]
	|_|  |___|[___| \  ||_|_|_|[___||_\_\\___.|_|  
	                [__/                           
	`)
	slog.Info("Server started at http://" + serverAddr)

	http.ListenAndServe(serverAddr, router)
}
