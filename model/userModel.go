package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string     `gorm:"unique;not null"`
	Email        string     `gorm:"unique;not null"`
	PasswordHash string     `gorm:"not null"`
	UserData     []UserData `gorm:"foreignKey:UserID"`
}

type UserData struct {
	gorm.Model
	UserID uint   `gorm:"index;not null"`
	Data   string `gorm:"type:jsonb;not null"`
}
