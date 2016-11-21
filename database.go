package main

import (
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	r "gopkg.in/dancannon/gorethink.v2"
	"log"
)

const (
	dataBase  = config.DataBase.DataBase
	dataTable = config.DataBase.Table
)

type (
	User struct {
		ID       int    `gorethink:"user_id"`
		NSFW     bool   `gorethink:"nsfw"`
		Language string `gorethink:"language"`
	}
)

func createUser(user *t.User) {
	resp, err := r.DB(dataBase).Table(dataTable).Filter(user.ID).Run(db)
	if err != nil {
		log.Println(err)
	}
	defer resp.Close()

	if resp.IsNil() {
		_, err := r.DB(dataBase).Table(dataTable).Insert(User{user.ID, false, "english"}).RunWrite(db)
		if err != nil {
			log.Println(err)
		}
	}
}

func checkNSFW(user *t.User) bool {
	createUser(user)
	resp, err := r.DB(dataBase).Table(dataTable).Filter(user.ID).Field("nsfw").Run(db)
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
	if _, err := r.DB(dataBase).Table(dataTable).Filter(user.ID).Update(User{NSFW: state}).RunWrite(db); err != nil {
		log.Println(err)
	}
}

func checkLanguage(user *t.User) string {
	createUser(user)
	resp, err := r.DB(dataBase).Table(dataTable).Filter(user.ID).Field("language").Run(db)
	if err != nil {
		log.Printf("[RethinkDB] %#v", err)
	}
	defer resp.Close()
	var lang string
	if err = resp.One(&lang); err != nil {
		log.Printf("[RethinkDB] %#v", err)
	}
	return lang
}

func changeLanguage(user *t.User, lang string) {
	createUser(user)
	if _, err := r.DB(dataBase).Table(dataTable).Filter(user.ID).Update(User{Language: lang}).RunWrite(db); err != nil {
		log.Println(err)
	}
}
