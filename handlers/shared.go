package handlers

import (
	"football-stat-goth/repos"
	"log/slog"
	"net/http"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

type HTTPRepoHandler func(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error

func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("Handler error", "err", err)
		}
	}
}

func MakeWithRepo(h HTTPRepoHandler, repo *repos.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r, repo); err != nil {
			slog.Error("Handler error", "err", err)
		}
	}
}
