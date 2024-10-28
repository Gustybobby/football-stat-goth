package plauth

import (
	"football-stat-goth/models"
	"football-stat-goth/repos"
	"net/http"
)

func Auth(r *http.Request, repo *repos.Repository) (*models.User, error) {
	token, err := GetSessionTokenFromCookie(r)
	if err != nil {
		return nil, err
	}

	session, err := ValidateSessionToken(token, repo)
	if err != nil {
		return nil, err
	}

	return &session.User, nil
}
