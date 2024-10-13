package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID            string     `gorm:"index:players_team_id_key;uniqueIndex:players_team_id_shirt_number_key;not null" json:"team_id"`
	ShirtNumber       uint       `gorm:"uniqueIndex:players_team_id_shirt_number_key;not null" json:"shirt_number"`
	ShirtName         string     `gorm:"not null" json:"shirt_name"`
	FirstName         string     `gorm:"not null" json:"first_name"`
	LastName          string     `gorm:"not null" json:"last_name"`
	DOB               *time.Time `json:"dob"`
	Nationality       *string    `json:"nationality"`
	PreferredPosition *string    `json:"preferred_position"`
}

func MigratePlayer(db *gorm.DB) error {
	err := db.AutoMigrate(&Player{})
	return err
}
