package models

import "gorm.io/gorm"

type Lineup struct {
	LineupID      uint    `gorm:"primaryKey;autoIncrement" json:"lineup_id"`
	ClubID        string  `gorm:"not null" json:"club_id"`
	Possession    float32 `gorm:"not null;default:0" json:"possession"`
	ShotsOnTarget uint    `gorm:"not null;default:0" json:"shots_on_target"`
	Shots         uint    `gorm:"not null;default:0" json:"shots"`
	Touches       uint    `gorm:"not null;default:0" json:"touches"`
	Passes        uint    `gorm:"not null;default:0" json:"passes"`
	Tackles       uint    `gorm:"not null;default:0" json:"tackles"`
	Clearances    uint    `gorm:"not null;default:0" json:"clearances"`
	Corners       uint    `gorm:"not null;default:0" json:"corners"`
	Offsides      uint    `gorm:"not null;default:0" json:"offsides"`
	FoulsConceded uint    `gorm:"not null;default:0" json:"fouls_conceded"`

	LineupPlayers []LineupPlayer `gorm:"foreignKey:LineupID;references:LineupID"`

	Club Club
}

func migrateLineup(db *gorm.DB) error {
	err := db.AutoMigrate(&Lineup{})
	return err
}
