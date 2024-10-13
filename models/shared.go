package models

import "gorm.io/gorm"

func MigrateSchema(db *gorm.DB) error {
	if err := MigrateTeam(db); err != nil {
		return err
	}
	return nil
}
