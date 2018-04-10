package database

import (
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func (db *DataBase) removeTag(list []models.Tag, name string) error {
	var err error
	for i := range list {
		if list[i].Tag == name {
			err = db.Delete(list[i]).Error
			break
		}
	}

	return err
}

func (db *DataBase) RemoveWhiteTag(usr *tg.User, tag string) error {
	user, err := db.GetUser(usr)
	if err != nil {
		return err
	}

	return db.removeTag(user.WhiteList, tag)
}

func (db *DataBase) RemoveBlackTag(usr *tg.User, tag string) error {
	user, err := db.GetUser(usr)
	if err != nil {
		return err
	}

	return db.removeTag(user.BlackList, tag)
}
