package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	MatchID      uint      `gorm:"primaryKey" json:"match_id"`
	HomeLineupID uint      `gorm:"not null;uniqueIndex" json:"home_lineup_id"`
	AwayLineupID uint      `gorm:"not null;uniqueIndex" json:"away_lineup_id"`
	Season       string    `gorm:"not null" json:"season"`
	HomeGoals    uint      `gorm:"not null;default:0" json:"home_goals"`
	AwayGoals    uint      `gorm:"not null;default:0" json:"away_goals"`
	Location     string    `gorm:"not null" json:"location"`
	DateTime     time.Time `gorm:"not null" json:"date_time"`
	MatchWeek    uint      `gorm:"not null" json:"match_week"`
	IsFinished   bool      `gorm:"not null;default:false" json:"is_finished"`

	HomeLineup Lineup `gorm:"foreignKey:HomeLineupID;references:LineupID"`
	AwayLineup Lineup `gorm:"foreignKey:AwayLineupID;references:LineupID"`
}

func migrateMatch(db *gorm.DB) error {
	err := db.AutoMigrate(&Match{})
	return err
}
