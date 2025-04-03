package repositories

import (
	"go-proxy-auth-service/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) FindByLogin(login string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("login = ?", login).First(&user).Error
	return &user, err
}
