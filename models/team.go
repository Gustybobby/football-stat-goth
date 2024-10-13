package models

import "gorm.io/gorm"

type Team struct {
	ID      string  `gorm:"primaryKey" json:"id"`
	Abv     *string `json:"abv"`
	Founded *string `json:"founded"`
	Owner   *string `json:"owner"`
}

func MigrateTeam(db *gorm.DB) error {
	err := db.AutoMigrate(&Team{})
	return err
}
