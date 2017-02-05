package main

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

var locales []string

func init() {
	if err := filepath.Walk("./i18n", func(path string, file os.FileInfo, err error) error {
		if !strings.HasPrefix(path, "i18n/source/") && strings.HasSuffix(path, ".all.json") {
			log.Ln("Load translation file", file.Name())
			i18n.MustLoadTranslationFile(path)
			locales = append(locales, strings.TrimSuffix(file.Name(), ".all.json"))
		}
		return nil
	}); err != nil {
		panic(err.Error())
	}
}
