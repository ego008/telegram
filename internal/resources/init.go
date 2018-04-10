package resources

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/errors"
	log "github.com/kirillDanshin/dlog"
	"github.com/spf13/viper"
)

var (
	Resources = make(map[string]*viper.Viper)

	Tags []string
)

func Initialize(pathToConfigs string) {
	filesMap := make(map[string]string)

	err := filepath.Walk(pathToConfigs, func(path string, file os.FileInfo, err error) error {
		log.Ln("Walk to", path)
		if !strings.HasSuffix(file.Name(), ".yaml") {
			return err
		}

		cfg := viper.New()
		cfg.AddConfigPath(path)
		cfg.SetConfigFile(file.Name())

		if err = cfg.ReadInConfig(); err != nil {
			return err
		}

		name := cfg.GetString("name")
		Resources[name] = cfg
		filesMap[name] = file.Name()

		return nil
	})
	errors.Check(err)

	for _, cfg := range Resources {
		tpl := cfg.GetString("template")
		if tpl == "" {
			continue
		}

		cfg.SetConfigFile(filesMap[tpl])
		if err = cfg.MergeInConfig(); err != nil {
			continue
		}
	}

	log.Ln("Resources after:")
	log.D(Resources)
}
