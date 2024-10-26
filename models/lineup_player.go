package models

import "gorm.io/gorm"

type LineupPlayer struct {
	PlayerNo       uint   `gorm:"primaryKey" json:"player_no"`
	ClubID         string `gorm:"primaryKey" json:"club_id"`
	LineupID       uint   `gorm:"primaryKey" json:"lineup_id"`
	PositionNumber uint   `gorm:"not null" json:"position_number"`
	Goals          uint   `gorm:"not null" json:"goals"`
	YellowCards    uint   `gorm:"not null" json:"yellow_cards"`
	RedCards       uint   `gorm:"not null" json:"red_cards"`
}

func migrateLineupPlayer(db *gorm.DB) error {
	err := db.AutoMigrate(&LineupPlayer{})
	return err
}
