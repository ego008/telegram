package main

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

type (
	user struct {
		Blacklist []string
		ID        int
		Language  string
		Ratings   ratings
		Resources map[string]bool
		Roles     roles
		Whitelist []string
	}

	roles struct {
		User, Patron, Manager, Admin bool
	}

	ratings struct {
		Safe, Questionable, Exlplicit bool
	}
)

const (
	blackList = "blacklist"
	whiteList = "whitelist"
)

var (
	db *buntdb.DB

	errNotFound = buntdb.ErrNotFound
)

func dbInit() {
	log.Ln("db:init")

	var err error
	db, err = buntdb.Open("hentai.db")
	errCheck(err)

	err = db.CreateIndex("users", "user:*", buntdb.IndexString)
	errCheck(err)

	select {}
}
