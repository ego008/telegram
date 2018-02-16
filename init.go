package main

import (
	"flag"

	"github.com/HentaiDB/HentaiDBot/internal/config"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	log "github.com/kirillDanshin/dlog"
)

const ver = `4.0 "Dark Dream"`

var (
	verHash, verTimeStamp string

	flagWebhook = flag.Bool("webhook", false, "activate getting updates via webhook")
)

func init() {
	log.Ln("Running", ver, "version...")
	flag.Parse()

	config.Initialize("./configs/config.yaml", *flagWebhook)
	resources.Initialize("./configs/resources")
	i18n.Initialize("./configs/translations")
	// go dbInit()
}
