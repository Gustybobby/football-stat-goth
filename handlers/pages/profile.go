package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
	"net/url"
)

func HandleProfilePage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Redirect(w, r, "/signin?redirectUrl="+url.QueryEscape("/profile"), http.StatusFound)
		return nil
	}

	return handlers.Render(w, r, views.Profile(user))
}
