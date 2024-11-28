package main

import (
	"embed"
	"football-stat-goth/handlers/api"
	"football-stat-goth/handlers/cmps"
	"football-stat-goth/handlers/pages"
	"football-stat-goth/handlers/plmiddleware"
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
	router.Use(plmiddleware.AuthMiddleware(repo))

	router.Handle("/*", http.FileServerFS(publicFS))

	pages.SetupPageRoutes(router, repo)

	router.Route("/api", func(r_api chi.Router) {
		api.SetupApiRoutes(r_api, repo)
	})

	router.Route("/cmps", func(r_cmps chi.Router) {
		cmps.SetupComponentRoutes(r_cmps, repo)
	})
}

func main() {
	if err := godotenv.Overload(); err != nil {
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
