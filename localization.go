package main

import (
	"github.com/nicksnyder/go-i18n/i18n"
)

var T, _ = i18n.Tfunc("en-us")

func init() {
	// Read localization
	i18n.MustLoadTranslationFile("i18n/en-us.all.json")
	i18n.MustLoadTranslationFile("i18n/ru-ru.all.json")
	i18n.MustLoadTranslationFile("i18n/zh-zh.all.json")
}
