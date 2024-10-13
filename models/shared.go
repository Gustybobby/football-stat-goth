package models

import "gorm.io/gorm"

func MigrateSchema(db *gorm.DB) error {
	if err := MigrateTeam(db); err != nil {
		return err
	}
	if err := MigratePlayer(db); err != nil {
		return err
	}
	return nil
}
