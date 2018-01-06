package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/kirillDanshin/dlog"
	config "github.com/olebedev/config"
)

var (
	cfg           *config.Config
	resources     = make(map[string]*config.Config)
	resourcesTags []string
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

	err = filepath.Walk("./resources", func(path string, file os.FileInfo, err error) error {
		log.Ln("Walk to", path)
		if !strings.HasSuffix(file.Name(), ".yaml") {
			return nil
		}

		res, err := config.ParseYamlFile(path)
		if err != nil {
			return err
		}

		resources[res.UString("name")] = res
		return nil
	})
	errCheck(err)

	log.Ln("Resources before:")
	log.D(resources)

	for res, conf := range resources {
		template := conf.UString("template")
		if template == "" {
			log.Ln("Resource:", res, "template: none")
			continue
		}

		resources[res], err = resources[template].Extend(conf)
		errCheck(err)

		resourcesTags = append(resourcesTags, res)
	}

	sort.Strings(resourcesTags)

	log.Ln("Resources after:")
	log.D(resources)
}
