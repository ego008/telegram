package i18n

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/HentaiDB/HentaiDBot/pkg/models"
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/yaml.v2"
)

type Bundle struct{ *i18n.Bundle }

var I18N *Bundle

func Open(path string) (*Bundle, error) {
	bundle := i18n.NewBundle(models.LanguageFallback)
	bundle.RegisterUnmarshalFunc("yaml", yaml.UnmarshalStrict)

	err := filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		log.Ln("Walk to", path)
		if !strings.HasSuffix(file.Name(), ".all.yaml") {
			return err
		}

		_, err = bundle.LoadMessageFile(path)
		return err
	})

	return &Bundle{bundle}, err
}
