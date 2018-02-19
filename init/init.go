package init

import (
	"flag"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/config"
	"github.com/HentaiDB/HentaiDBot/internal/db"
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
	flag.Parse()
	log.Ln("Running", ver, "version...")

	config.Initialize("./configs/config.yaml", *flagWebhook)
	resources.Initialize("./configs/resources")
	i18n.Initialize("./configs/translations")
	db.Initialize()
	bot.Initialize()
}
