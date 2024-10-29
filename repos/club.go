package repos

import "football-stat-goth/models"

func FindClubsWithNameAsc(repo *Repository) ([]models.Club, error) {
	var clubs []models.Club
	results := repo.DB.Select("ClubID", "Name", "Logo").Order("name ASC").Find(&clubs)
	if results.Error != nil {
		return nil, results.Error
	}
	return clubs, nil
}

func FindClubs(repo *Repository) ([]models.Club, error) {
	var clubs []models.Club
	results := repo.DB.Find(&clubs)
	if results.Error != nil {
		return nil, results.Error
	}
	return clubs, nil
}
