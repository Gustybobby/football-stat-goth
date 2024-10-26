package models

import (
	"time"

	"gorm.io/gorm"
)

type PlayerPosition string

const (
	GK  PlayerPosition = "GK"
	DEF PlayerPosition = "DEF"
	MFD PlayerPosition = "MFD"
	FWD PlayerPosition = "FWD"
	SUB PlayerPosition = "SUB"
)

type Player struct {
	PlayerNo    uint           `gorm:"primaryKey" json:"player_no"`
	ClubID      string         `gorm:"primaryKey" json:"club_id"`
	FirstName   string         `gorm:"not null" json:"first_name"`
	LastName    string         `gorm:"not null" json:"last_name"`
	DOB         time.Time      `gorm:"not null" json:"dob"`
	Height      uint           `gorm:"not null" json:"height"`
	Nationality string         `gorm:"not null" json:"nationality"`
	Position    PlayerPosition `gorm:"not null;type:player_position" json:"position"`

	LineupPlayers []LineupPlayer `gorm:"foreignKey:PlayerNo,ClubID;references:PlayerNo,ClubID"`
}

func migratePlayer(db *gorm.DB) error {
	err := db.AutoMigrate(&Player{})
	return err
}
