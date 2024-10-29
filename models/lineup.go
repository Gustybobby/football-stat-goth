package models

import "gorm.io/gorm"

type Lineup struct {
	LineupID      uint    `gorm:"primaryKey;autoIncrement" json:"lineup_id"`
	ClubID        string  `gorm:"not null" json:"club_id"`
	Possession    float32 `gorm:"not null" json:"possession"`
	ShotsOnTarget uint    `gorm:"not null" json:"shots_on_target"`
	Shots         uint    `gorm:"not null" json:"shots"`
	Touches       uint    `gorm:"not null" json:"touches"`
	Passes        uint    `gorm:"not null" json:"passes"`
	Tackles       uint    `gorm:"not null" json:"tackles"`
	Clearances    uint    `gorm:"not null" json:"clearances"`
	Corners       uint    `gorm:"not null" json:"corners"`
	Offsides      uint    `gorm:"not null" json:"offsides"`
	FoulsConceded uint    `gorm:"not null" json:"fouls_conceded"`

	LineupPlayers []LineupPlayer `gorm:"foreignKey:LineupID;references:LineupID"`

	Club Club
}

func migrateLineup(db *gorm.DB) error {
	err := db.AutoMigrate(&Lineup{})
	return err
}
