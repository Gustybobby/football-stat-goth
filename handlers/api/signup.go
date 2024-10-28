package api

import (
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"net/http"
	"os"
)

func HandleSignup(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user, err := plauth.CreateUser(r.FormValue("username"), r.FormValue("password"), r.FormValue("first_name"), r.FormValue("last_name"), repo)
	if err != nil {
		return err
	}

	token, err := plauth.GenerateSessionToken()
	if err != nil {
		return err
	}

	session, err := plauth.CreateSession(token, user.Username, repo)
	if err != nil {
		return err
	}

	plauth.SetSessionTokenCookie(w, token, session.ExpiresAt, os.Getenv("ENV") == "production")

	w.Header().Add("Hx-Redirect", "/")
	return nil
}
