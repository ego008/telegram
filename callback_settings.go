package main

import (
	"strings"

	// log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func callbackToSettings(usr *user, call *tg.CallbackQuery) {
	T, err := langSwitch(usr.Language, call.From.LanguageCode)
	errCheck(err)

	var activeRes []string
	for k, v := range usr.Resources {
		if v && resources[k] != nil {
			activeRes = append(activeRes, resources[k].UString("title"))
		}
	}

	ratings, err := usr.getRatingsStatus()
	errCheck(err)

	text := T("message_settings", map[string]interface{}{
		"Language":  languageNames[usr.Language],
		"Resources": strings.Join(activeRes, "`, `"),
		"Ratings":   ratings,
		"Blacklist": strings.Join(usr.Blacklist, "`, `"),
		"Whitelist": strings.Join(usr.Whitelist, "`, `"),
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = getSettingsMenuKeyboard(usr)

	_, err = bot.EditMessageText(editText)
	errCheck(err)
}

func getSettingsMenuKeyboard(usr *user) *tg.InlineKeyboardMarkup {
	T, err := langSwitch(usr.Language)
	errCheck(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_language", map[string]interface{}{
					"Flag": T("language_flag"),
				}),
				"to:languages",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_resources"), "to:resources",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_ratings"), "to:ratings",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				"Type filters", "to:types",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_blacklist"), "to:blacklist",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_whitelist"), "to:whitelist",
			),
		),
	)
}
