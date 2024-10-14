package models

import "gorm.io/gorm"

type Team struct {
	ID      string  `gorm:"primaryKey" json:"id"`
	Name    string  `gorm:"not null" json:"name"`
	Founded uint    `gorm:"not null" json:"founded"`
	Owner   *string `json:"owner"`

	Players []Player `gorm:"foreignKey:TeamID;references:ID"`
}

func migrateTeam(db *gorm.DB) error {
	err := db.AutoMigrate(&Team{})
	return err
}
