package repositories

import (
	"go-proxy-auth-service/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) FindUserWithSession(login string, sessionId string) (*models.User, error) {
	var user models.User
	err := r.DB.Model(&models.Session{}).Select("*").
		Joins("JOIN users ON users.login = sessions.login").
		Where("users.login = ? AND sessions.sessionId = ?", login, sessionId).First(&user).Error
	return &user, err
}
