package repos

import "football-stat-goth/models"

func FindClubsByNameAsc(repo *Repository) ([]models.Club, error) {
	var clubs []models.Club
	results := repo.DB.Order("name ASC").Find(&clubs)
	if results.Error != nil {
		return nil, results.Error
	}
	return clubs, nil
}
