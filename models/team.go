package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID      string    `gorm:"primaryKey" json:"id"`
	Name    string    `gorm:"not null" json:"name"`
	Founded time.Time `gorm:"not null" json:"founded"`
	Owner   *string   `json:"owner"`

	Players []Player `gorm:"foreignKey:TeamID;references:ID"`
}

func MigrateTeam(db *gorm.DB) error {
	err := db.AutoMigrate(&Team{})
	return err
}
