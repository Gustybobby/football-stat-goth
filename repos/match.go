package repos

import (
	"football-stat-goth/models"
	"time"
)

func FindFixtureMatches(repo *Repository) ([]models.Match, error) {
	var matches []models.Match
	var earliestTimeString string = time.Now().Add(time.Duration(-3) * time.Hour).UTC().Format(time.RFC3339)
	results := repo.DB.Model(&models.Match{}).Preload(
		"HomeLineup.Club").Preload(
		"AwayLineup.Club").Where(
		"date_time >= ?", earliestTimeString).Order(
		"date_time ASC").Find(&matches)
	if results.Error != nil {
		return nil, results.Error
	}
	return matches, nil
}
