package repos

import (
	"football-stat-goth/models"
)

func FindFixtureMatches(repo *Repository) ([]models.Match, error) {
	var matches []models.Match
	results := repo.DB.Model(&models.Match{}).Preload(
		"HomeLineup.Club").Preload(
		"AwayLineup.Club").Where(
		"is_finished = false").Order(
		"date_time ASC").Find(&matches)
	if results.Error != nil {
		return nil, results.Error
	}
	return matches, nil
}

func FindFixtureMatchesByClubID(clubID string, repo *Repository) ([]models.Match, error) {
	var matches []models.Match
	results := repo.DB.Model(&models.Match{}).Preload(
		"HomeLineup.Club").Preload(
		"AwayLineup.Club").Joins("INNER JOIN lineups ON home_lineup_id = lineups.lineup_id AND lineups.club_id = ?", clubID).Where(
		"is_finished = false").Order(
		"date_time DESC").Find(&matches)
	if results.Error != nil {
		return nil, results.Error
	}
	return matches, nil
}

func FindResultMatchesByClubID(clubID string, repo *Repository) ([]models.Match, error) {
	var matches []models.Match
	results := repo.DB.Model(&models.Match{}).Preload(
		"HomeLineup.Club").Preload(
		"AwayLineup.Club").Joins("INNER JOIN lineups ON home_lineup_id = lineups.lineup_id AND lineups.club_id = ?", clubID).Where(
		"is_finished = true").Order(
		"date_time DESC").Find(&matches)
	if results.Error != nil {
		return nil, results.Error
	}
	return matches, nil
}
