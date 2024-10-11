package handlers

import (
	"log/slog"
	"net/http"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("Handler error", "err", err)
		}
	}
}
