package main

import (
	"football-stat-goth/handlers"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()

	router.Get("/", handlers.Make(handlers.HandleHelloWorld))

	serverAddr := os.Getenv("SERVER_ADDR")

	slog.Info("Server started", "serverAddr", serverAddr)
	http.ListenAndServe(serverAddr, router)
}
