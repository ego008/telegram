package main

import (
	log "github.com/kirillDanshin/dlog"
	config "github.com/olebedev/config"
)

var cfg *config.Config

func cfgInit() {
	var err error
	log.Ln("Parse yaml config...")
	cfg, err = config.ParseYamlFile("config.yaml")
	errCheck(err)

	log.Ln("Check telegram token...")
	_, err = cfg.String("telegram.token")
	errCheck(err)

	if *flagWebhook {
		log.Ln("Check telegram webhook set...")
		_, err = cfg.String("telegram.webhook.set")
		errCheck(err)

		log.Ln("Check telegram webhook listen...")
		_, err = cfg.String("telegram.webhook.listen")
		errCheck(err)

		log.Ln("Check telegram webhook serve...")
		_, err = cfg.String("telegram.webhook.serve")
		errCheck(err)
	}
}
