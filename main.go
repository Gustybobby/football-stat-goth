// THIS IS TN BRANCH, CREATED JUST TO PUSH INTO GITHUB...
package main

import (
	"embed"
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/api"
	"football-stat-goth/handlers/pages"
	"football-stat-goth/models"
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
	router.Get("/signup", handlers.Make(pages.HandleSignupPage, repo))
	router.Get("/clubs", handlers.Make(pages.HandleClubsPage, repo))

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

	db, err := repos.NewConnection(config)
	if err != nil {
		log.Fatal("could not connect to database")
	}

	if os.Getenv("DB_AUTOMIGRATE") == "true" {
		if err := models.MigrateSchema(db); err != nil {
			log.Fatal(err)
		}
	}

	repo := &repos.Repository{
		DB: db,
	}

	SetupRoutes(router, repo)

	serverAddr := os.Getenv("SERVER_ADDR")

	slog.Info("Server started", "serverAddr", serverAddr)
	http.ListenAndServe(serverAddr, router)
}
