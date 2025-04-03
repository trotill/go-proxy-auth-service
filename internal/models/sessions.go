package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	id           int32     `gorm:"unique;not null"`
	login        string    `gorm:"not null"`
	sessionId    string    `gorm:"not null"`
	refreshToken string    `gorm:"not null"`
	createdAt    time.Time `gorm:"not null"`
	updatedAt    time.Time `gorm:"not null"`
}
