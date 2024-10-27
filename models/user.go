package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username     string    `gorm:"primaryKey" json:"username"`
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`

	Sessions []Session `gorm:"foreignKey:Username;references:Username"`
}

func migrateUser(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
