package api

import (
	"crypto/subtle"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"net/http"
)

func HandleUpdatePassword(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	userPasswordHash, err := repo.Queries.FindPasswordHashByUsername(repo.Ctx, user.Username)
	if err != nil {
		return err
	}

	currentPasswordHash, err := plauth.HashPassword(r.FormValue("current"))
	if err != nil {
		return err
	}

	newPasswordHash, err := plauth.HashPassword(r.FormValue("new"))
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare([]byte(currentPasswordHash), []byte(newPasswordHash)) == 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	if subtle.ConstantTimeCompare([]byte(userPasswordHash), []byte(currentPasswordHash)) == 1 {
		repo.Queries.UpdatePasswordByUsername(repo.Ctx, queries.UpdatePasswordByUsernameParams{
			Username:     user.Username,
			PasswordHash: newPasswordHash,
		})
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return nil
}
