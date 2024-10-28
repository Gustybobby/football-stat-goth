package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	USER  UserRole = "USER"
	ADMIN UserRole = "ADMIN"
)

type User struct {
	Username     string    `gorm:"primaryKey" json:"username"`
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	Role         UserRole  `gorm:"not null;type:user_role;default:'USER'" json:"role"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`

	Sessions []Session `gorm:"foreignKey:Username;references:Username"`
}

func migrateUser(db *gorm.DB) error {
	err1 := db.Exec("CREATE TYPE user_role AS ENUM ('USER','ADMIN');")
	if err1 != nil {
		log.Default().Println("Type user_role already exist")
	}
	err := db.AutoMigrate(&User{})
	return err
}
