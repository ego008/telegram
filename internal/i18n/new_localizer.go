package i18n

import "github.com/nicksnyder/go-i18n/v2/i18n"

type LocalizeConfig = i18n.LocalizeConfig

func (bundle *Bundle) NewLocalizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle.Bundle, langs...)
}
