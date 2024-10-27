package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	Token     string    `gorm:"primaryKey" json:"token"`
	Username  string    `gorm:"not null" json:"username"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}

func migrateSession(db *gorm.DB) error {
	err := db.AutoMigrate(&Session{})
	return err
}
