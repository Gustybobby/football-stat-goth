package plmiddleware

import (
	"context"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func AuthMiddleware(repo *repos.Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.String(), "/public") {
				next.ServeHTTP(w, r)
				return
			}
			user, err := plauth.Auth(w, r, repo)
			if err != nil {
				return
			}
			ctx := context.WithValue(r.Context(), plauth.ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := plauth.GetContextUser(r)
		if user == nil || user.Role != queries.UserRoleADMIN {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthSessionUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := plauth.GetContextUser(r)
		param_username := chi.URLParam(r, "username")
		if user.Username != param_username && user.Role != queries.UserRoleADMIN {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
