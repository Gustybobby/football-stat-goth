package repos

import "football-stat-goth/models"

func FindFixtureMatches(repo *Repository) ([]models.Match, error) {
	var matches []models.Match
	results := repo.DB.Model(&models.Match{}).Preload("HomeLineup.Club").Preload("HomeLineup").Preload("AwayLineup.Club").Preload("AwayLineup").Find(&matches)
	if results.Error != nil {
		return nil, results.Error
	}
	return matches, nil
}
