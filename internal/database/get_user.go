package database

import (
	"time"

	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func (db *DataBase) GetUser(usr *tg.User) (*models.User, error) {
	var err error
	var user, newUser models.User
	user.Model.ID = usr.ID

	newUser.Locale = usr.LanguageCode
	newUser.Model = models.Model{
		ID:        usr.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err = db.FirstOrCreate(&user, &newUser).Error; err != nil {
		return &user, err
	}

	return &user, err
}
