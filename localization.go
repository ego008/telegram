package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	i18n "github.com/nicksnyder/go-i18n/i18n"
)

var locales []string

func langInit() {
	if err := filepath.Walk(*flagLocale, func(path string, file os.FileInfo, err error) error {
		if !strings.HasPrefix(path, "i18n/source/") && strings.HasSuffix(path, ".all.json") {
			log.Println("Load translation file", file.Name())
			i18n.MustLoadTranslationFile(path)
			locales = append(locales, strings.TrimSuffix(file.Name(), ".all.json"))
		}
		return nil
	}); err != nil {
		log.Fatalln(err.Error())
	}
}
