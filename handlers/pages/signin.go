package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
	"net/url"
)

func HandleSigninPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	if user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	redirect_url := r.URL.Query().Get("redirectUrl")
	if redirect_url == "" {
		redirect_url = "/"
	}

	return handlers.Render(w, r, views.Signin(user, url.QueryEscape(redirect_url)))
}
