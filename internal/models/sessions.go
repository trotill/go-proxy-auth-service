package models

import (
	"time"
)

type Session struct {
	Id           int32     `gorm:"unique;not null"`
	Login        string    `gorm:"not null;index"`
	SessionId    string    `gorm:"not null;column:sessionId"`
	RefreshToken string    `gorm:"not null;column:refreshToken"`
	CreatedAt    time.Time `gorm:"not null;column:createdAt"`
	UpdatedAt    time.Time `gorm:"not null;column:updatedAt"`
}
