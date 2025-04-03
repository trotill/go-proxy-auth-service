package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	login     string    `gorm:"unique;not null"`
	role      string    `gorm:"not null"`
	password  string    `gorm:"not null"`
	firstName string    `gorm:"not null"`
	lastName  string    `gorm:"not null"`
	email     string    `gorm:"not null"`
	locked    uint32    `gorm:"not null"`
	createdAt time.Time `gorm:"not null"`
	updatedAt time.Time `gorm:"not null"`
}
