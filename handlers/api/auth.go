package api

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

func HandleSignup(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user, err := plauth.CreateUser(r.FormValue("username"), r.FormValue("password"), r.FormValue("first_name"), r.FormValue("last_name"), repo.Queries, repo.Ctx)
	if err != nil {
		slog.Error(err.Error())
		return handlers.Render(w, r, views.SignupForm("Username already existed"))
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

	w.Header().Add("HX-Redirect", "/")
	return nil
}

func HandleSignin(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	redirect_url, err := url.QueryUnescape(r.URL.Query().Get("redirectUrl"))
	if err != nil {
		return err
	}
	if redirect_url == "" {
		redirect_url = "/"
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	passwordHash, err := plauth.FindPasswordHash(username, repo.Queries, repo.Ctx)
	if err != nil {
		slog.Error("username not found")
		return handlers.Render(w, r, views.SigninForm(url.QueryEscape(redirect_url), "Invalid username or password"))
	}

	valid, err := plauth.VerifyPassword(password, passwordHash)
	if err != nil {
		return err
	}

	if !valid {
		slog.Error("invalid password")
		return handlers.Render(w, r, views.SigninForm(url.QueryEscape(redirect_url), "Invalid username or password"))
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

	w.Header().Add("HX-Redirect", redirect_url)
	return nil
}

func HandleSignout(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	plauth.DeleteSessionTokenCookie(w, os.Getenv("ENV") == "production")
	w.Header().Add("HX-Redirect", "/")
	return nil
}

func HandleUpdatePassword(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	currentPassword := r.FormValue("current")
	newPassword := r.FormValue("new")
	confirmNewPassword := r.FormValue("confirm_new")

	if newPassword != confirmNewPassword {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	err := plauth.UpdatePassword(user.Username, currentPassword, newPassword, repo.Queries, repo.Ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return err
	}

	w.Header().Add("HX-Refresh", "true")
	return nil
}
