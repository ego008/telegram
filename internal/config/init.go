package config

import (
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	log "github.com/kirillDanshin/dlog"
	config "github.com/olebedev/config"
)

var Config *config.Config

func Initialize(pathToConfig string, webhookMode bool) {
	var err error
	log.Ln("Parse yaml config...")
	Config, err = config.ParseYamlFile(pathToConfig)
	errors.Check(err)

	log.Ln("Check telegram token...")
	_, err = Config.String("telegram.token")
	errors.Check(err)

	if webhookMode {
		log.Ln("Check telegram webhook set...")
		_, err = Config.String("telegram.webhook.set")
		errors.Check(err)

		log.Ln("Check telegram webhook listen...")
		_, err = Config.String("telegram.webhook.listen")
		errors.Check(err)

		log.Ln("Check telegram webhook serve...")
		_, err = Config.String("telegram.webhook.serve")
		errors.Check(err)
	}
}
