package main

import (
	"embed"
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/api"
	"football-stat-goth/handlers/pages"
	"football-stat-goth/repos"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
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
	router.Get("/signup", handlers.Make(pages.HandleSignupPage, repo))

	router.Post("/api/signup", handlers.Make(api.HandleSignup, repo))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()

	config := &repos.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	repo := repos.New(config)

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
