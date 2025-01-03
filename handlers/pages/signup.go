package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views"
	"net/http"
)

func HandleSignupPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	if user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	return handlers.Render(w, r, views.Signup(user))
}
