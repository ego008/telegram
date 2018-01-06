package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

var languageNames = make(map[string]string)
var languageTags []string

func langInit() {
	err := filepath.Walk("./i18n", func(path string, file os.FileInfo, err error) error {
		log.Ln("Walk to", path)
		if !strings.HasSuffix(file.Name(), ".all.yaml") {
			return nil
		}

		i18n.MustLoadTranslationFile(path)
		return nil
	})
	errCheck(err)

	languageTags = i18n.LanguageTags()
	for _, tag := range languageTags {
		T, err := langSwitch(tag)
		errCheck(err)

		languageNames[tag] = fmt.Sprint(
			T("language_flag"), " ", strings.Title(T("language_name")),
		)
		log.Ln("Tag", tag, ":", languageNames[tag])
	}

	sort.Strings(languageTags)
}
