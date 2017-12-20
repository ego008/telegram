package main

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

const langDefault = "en"

func langSwitch(langCodes ...string) (T i18n.TranslateFunc, err error) {
	log.Ln("Trying set localization")
	langCodes = append(langCodes, langDefault)
	T, err = i18n.Tfunc(langCodes[0], langCodes[1:]...)
	return
}
