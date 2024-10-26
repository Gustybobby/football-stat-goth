package models

import "gorm.io/gorm"

func MigrateSchema(db *gorm.DB) error {
	if err := migrateClub(db); err != nil {
		return err
	}
	if err := migratePlayer(db); err != nil {
		return err
	}
	if err := migrateLineup(db); err != nil {
		return err
	}
	if err := migrateLineupPlayer(db); err != nil {
		return err
	}
	if err := migrateMatch(db); err != nil {
		return err
	}
	return nil
}
