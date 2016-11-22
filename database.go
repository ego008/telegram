package main

import (
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/i18n"
	r "gopkg.in/dancannon/gorethink.v2"
	"log"
)

type (
	User struct {
		ID       int    `gorethink:"user_id"`
		NSFW     bool   `gorethink:"nsfw"`
		Language string `gorethink:"language"`
	}
)

func createUser(user *t.User) i18n.TranslateFunc {
	resp, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Filter(map[string]interface{}{"user_id": user.ID}).Run(db)
	if err != nil {
		log.Println(err)
	}
	defer resp.Close()

	var locale i18n.TranslateFunc

	if resp.IsNil() {
		_, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Insert(User{user.ID, false, "en-us"}).RunWrite(db)
		if err != nil {
			log.Println(err)
		}
		locale, _ = i18n.Tfunc("en-us")
	}
	return locale
}

func checkNSFW(user *t.User) bool {
	createUser(user)
	resp, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Filter(map[string]interface{}{"user_id": user.ID}).Field("nsfw").Run(db)
	if err != nil {
		log.Printf("[RethinkDB] %#v", err)
	}
	defer resp.Close()
	var nsfw bool
	if err = resp.One(&nsfw); err != nil {
		log.Printf("[RethinkDB] %#v", err)
	}
	return nsfw
}

func switchNSFW(user *t.User, state bool) {
	createUser(user)
	if _, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Filter(map[string]interface{}{"user_id": user.ID}).Update(map[string]interface{}{"nsfw": state}).RunWrite(db); err != nil {
		log.Println(err)
	}
}

func checkLanguage(user *t.User) i18n.TranslateFunc {
	createUser(user)
	resp, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Filter(map[string]interface{}{"user_id": user.ID}).Field("language").Run(db)
	if err != nil {
		log.Printf("[RethinkDB] %#v", err)
	}
	defer resp.Close()
	var lang string
	if err = resp.One(&lang); err != nil {
		log.Printf("[RethinkDB] %#v", err)
	}
	locale, _ := i18n.Tfunc(lang)
	return locale
}

func changeLanguage(user *t.User, lang string) i18n.TranslateFunc {
	createUser(user)
	if _, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Filter(map[string]interface{}{"user_id": user.ID}).Update(map[string]interface{}{"language": lang}).RunWrite(db); err != nil {
		log.Println(err)
	}
	locale, _ := i18n.Tfunc(lang)
	return locale
}
