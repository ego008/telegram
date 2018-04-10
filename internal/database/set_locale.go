package database

import tg "github.com/toby3d/telegram"

func (db *DataBase) SetLocale(usr *tg.User, locale string) error {
	user, err := db.GetUser(usr)
	if err != nil {
		return err
	}

	return db.Model(user).Set("locale", locale).Error
}
