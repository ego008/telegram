package init

import (
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/config"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	log "github.com/kirillDanshin/dlog"
)

func init() {
	log.Ln("Running", models.Version, "version...")

	var err error
	config.Config, err = config.Open("./configs/config.yaml")
	errors.Check(err)

	resources.Initialize("./configs/resources")

	i18n.I18N, err = i18n.Open("./configs/translations")
	errors.Check(err)

	database.DB, err = database.Open("./hentai.db")
	errors.Check(err)

	bot.Bot, err = bot.New(config.Config.GetString("telegram.token"))
	errors.Check(err)
}
