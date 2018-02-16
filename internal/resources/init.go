package resources

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/errors"
	log "github.com/kirillDanshin/dlog"
	"github.com/olebedev/config"
)

var (
	Resources = make(map[string]*config.Config)

	Tags []string
)

func Initialize(pathToConfigs string) {
	err := filepath.Walk(pathToConfigs, func(path string, file os.FileInfo, err error) error {
		log.Ln("Walk to", path)
		if !strings.HasSuffix(file.Name(), ".yaml") {
			return nil
		}

		res, err := config.ParseYamlFile(path)
		if err != nil {
			return err
		}

		Resources[res.UString("name")] = res
		return nil
	})
	errors.Check(err)

	log.Ln("Resources before:")
	log.D(Resources)

	for res, conf := range Resources {
		template := conf.UString("template")
		if template == "" {
			log.Ln("Resource:", res, "template: none")
			continue
		}

		Resources[res], err = Resources[template].Extend(conf)
		errors.Check(err)

		Tags = append(Tags, res)
	}

	sort.Strings(Tags)

	log.Ln("Resources after:")
	log.D(Resources)
}
