package main

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	config "github.com/olebedev/config"
)

var (
	cfg       *config.Config
	resources = make(map[string]map[string]interface{})
)

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

	for i := 0; i == i; i++ {
		res, err := cfg.Map(fmt.Sprint("resources.", i))
		if err != nil {
			break
		}

		name := res["name"].(string)
		log.Ln("Getted", name, "resource config")
		resources[name] = res
	}
}
