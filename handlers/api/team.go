package api

import (
	"encoding/json"
	"fmt"
	"football-stat-goth/handlers"
	"football-stat-goth/models"
	"football-stat-goth/repos"
	"net/http"
)

func HandleCreateTeam(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	if err := handlers.AuthAPIKey(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return err
	}

	var team models.Team

	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	fmt.Printf("Team: %+v", team)

	result := repo.DB.Create(&team)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return result.Error
	}

	response, err := json.Marshal(&team)
	if err != nil {
		return err
	}

	w.Write(response)
	return nil
}
