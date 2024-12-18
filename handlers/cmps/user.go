package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/views/components/profile_components"
	"net/http"
)

func HandleChangePasswordForm(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	return handlers.Render(w, r, profile_components.Password(user))
}

func HandleEditProfileForm(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	return handlers.Render(w, r, profile_components.ProfileForm(user))
}
