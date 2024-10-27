package plauth

import (
	"football-stat-goth/models"
	"football-stat-goth/repos"
)

func CreateUser(username string, password string, firstName string, lastName string, repo *repos.Repository) (*models.User, error) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	var user = models.User{Username: username, PasswordHash: passwordHash, FirstName: firstName, LastName: lastName}

	results := repo.DB.Create(&user)
	if results.Error != nil {
		return nil, results.Error
	}

	return &user, nil
}
