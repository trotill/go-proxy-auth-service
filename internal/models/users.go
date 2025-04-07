package models

import (
	"time"
)

type User struct {
	Login     string    `gorm:"unique;not null"`
	Role      string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	Locked    uint32    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	//Sessions  []Session `gorm:"foreignKey:Login;references:Login"`
}
