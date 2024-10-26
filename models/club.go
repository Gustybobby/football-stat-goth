package models

import "gorm.io/gorm"

type Club struct {
	ClubID string `gorm:"primaryKey" json:"club_id"`
	Name   string `gorm:"not null" json:"name"`
	Logo   string `json:"logo"`

	Players []Player `gorm:"foreignKey:ClubID;references:ClubID"`
	Lineups []Lineup `gorm:"foreignKey:ClubID;references:ClubID"`
}

func migrateClub(db *gorm.DB) error {
	err := db.AutoMigrate(&Club{})
	return err
}
