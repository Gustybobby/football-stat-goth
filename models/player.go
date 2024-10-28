package models

import (
	"log"
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
	PlayerID    uint           `gorm:"primaryKey" json:"player_id"`
	ClubID      string         `json:"club_id"`
	FirstName   string         `gorm:"not null" json:"first_name"`
	LastName    string         `gorm:"not null" json:"last_name"`
	DOB         time.Time      `gorm:"not null" json:"dob"`
	Height      uint           `gorm:"not null" json:"height"`
	Nationality string         `gorm:"not null" json:"nationality"`
	Position    PlayerPosition `gorm:"not null;type:player_position" json:"position"`

	LineupPlayers []LineupPlayer `gorm:"foreignKey:PlayerID;references:PlayerID"`
}

func migratePlayer(db *gorm.DB) error {
	err1 := db.Exec("CREATE TYPE player_position AS ENUM ('GK','DEF','MFD','FWD','SUB');")
	if err1 != nil {
		log.Default().Println("Type player_position already exist")
	}
	err2 := db.AutoMigrate(&Player{})
	return err2
}
