package handlers

import (
	"errors"
	"football-stat-goth/repos"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error

func Make(h HTTPHandler, repo *repos.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r, repo); err != nil {
			slog.Error("Handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}

func AuthAPIKey(r *http.Request) error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if r.Header.Get("apikey") != os.Getenv("API_KEY") {
		return errors.New("invalid api key")
	}
	return nil
}
