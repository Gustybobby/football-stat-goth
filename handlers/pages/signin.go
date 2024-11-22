package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
)

func HandleSigninPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user, err := plauth.Auth(w, r, repo)
	if err != nil {
		return err
	}

	if user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	return handlers.Render(w, r, views.Signin(user))
}
