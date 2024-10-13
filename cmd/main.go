package main

import (
	"football-stat-goth/handlers"
	"football-stat-goth/models"
	"football-stat-goth/repos"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func SetupRoutes(router *chi.Mux, repo *repos.Repository) {
	router.Get("/", handlers.Make(handlers.HandleHelloWorld))
	router.Post("/teams", handlers.MakeWithRepo(handlers.HandleCreateTeam, repo))
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
