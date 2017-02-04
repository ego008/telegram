package main

import (
	"log"

	"github.com/boltdb/bolt"
	// tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var db *bolt.DB

func init() {
	var err error
	go func() {
		db, err = bolt.Open("hentai.db", 0600, nil)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer db.Close()

		select {}
	}()
}

/*
func createUser(user *tg.User) i18n.TranslateFunc {
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

func checkNSFW(user *tg.User) bool {
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

func switchNSFW(user *tg.User, state bool) {
	createUser(user)
	if _, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Filter(map[string]interface{}{"user_id": user.ID}).Update(map[string]interface{}{"nsfw": state}).RunWrite(db); err != nil {
		log.Println(err)
	}
}

func checkLanguage(user *tg.User) i18n.TranslateFunc {
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

func changeLanguage(user *tg.User, lang string) i18n.TranslateFunc {
	createUser(user)
	if _, err := r.DB(config.DataBase.DataBase).Table(config.DataBase.Table).Filter(map[string]interface{}{"user_id": user.ID}).Update(map[string]interface{}{"language": lang}).RunWrite(db); err != nil {
		log.Println(err)
	}
	locale, _ := i18n.Tfunc(lang)
	return locale
}
*/
