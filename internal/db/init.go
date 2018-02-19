package db

import (
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

const (
	blackList = "blacklist"
	whiteList = "whitelist"
)

var DB *buntdb.DB

func Initialize() {
	log.Ln("db:init")

	go func() {
		var err error
		DB, err = buntdb.Open("./hentai.db")
		errors.Check(err)

		select {}
	}()
}
