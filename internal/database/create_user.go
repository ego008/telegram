package database

import (
	"time"

	"github.com/HentaiDB/HentaiDBot/pkg/models"
)

func (db *DataBase) CreateUser(user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := db.Create(user).Error
	return user, err
}
