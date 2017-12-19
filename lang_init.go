package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

var langTags = make(map[string]string)

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

	tags := i18n.LanguageTags()
	for _, tag := range tags {
		T, err := langSwitch(tag)
		errCheck(err)

		langTags[tag] = fmt.Sprint(
			T("language_flag"), " ", strings.Title(T("language_name")),
		)
		log.Ln("Tag", tag, ":", langTags[tag])
	}
}
