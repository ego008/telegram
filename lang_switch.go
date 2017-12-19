package main

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

const langDefault = "en"

func langSwitch(langCode string) (T i18n.TranslateFunc, err error) {
	log.Ln("Trying set", langCode, "localization")
	T, err = i18n.Tfunc(langCode)
	if err != nil {
		log.Ln("Unsupported language, set", langDefault, "language as default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}
	return
}
