package api

import (
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"log/slog"
	"net/http"
	"os"
)

func HandleSignup(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user, err := plauth.CreateUser(r.FormValue("username"), r.FormValue("password"), r.FormValue("first_name"), r.FormValue("last_name"), repo.Queries, repo.Ctx)
	if err != nil {
		return err
	}

	token, err := plauth.GenerateSessionToken()
	if err != nil {
		return err
	}

	session, err := plauth.CreateSession(token, user.Username, repo.Queries, repo.Ctx)
	if err != nil {
		return err
	}

	plauth.SetSessionTokenCookie(w, token, session.ExpiresAt.Time, os.Getenv("ENV") == "production")

	w.Header().Add("Hx-Redirect", "/")
	return nil
}

func HandleSignin(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	passwordHash, err := plauth.FindPasswordHash(username, repo.Queries, repo.Ctx)
	if err != nil {
		slog.Error("username not found")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return err
	}

	valid, err := plauth.VerifyPassword(password, passwordHash)
	if err != nil {
		return err
	}

	if !valid {
		slog.Error("invalid password")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	token, err := plauth.GenerateSessionToken()
	if err != nil {
		return err
	}

	session, err := plauth.CreateSession(token, username, repo.Queries, repo.Ctx)
	if err != nil {
		return err
	}
	slog.Info("created new session for user: " + username)

	plauth.SetSessionTokenCookie(w, token, session.ExpiresAt.Time, os.Getenv("ENV") == "production")

	w.Header().Add("Hx-Redirect", "/")
	return nil
}
