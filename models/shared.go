package models

import "gorm.io/gorm"

func MigrateSchema(db *gorm.DB) error {
	if err := migrateTeam(db); err != nil {
		return err
	}
	if err := migratePlayer(db); err != nil {
		return err
	}
	return nil
}
