package main

import (
	"fmt"
	"strings"

	// log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func callbackSwitchLanguage(usr *user, call *tg.CallbackQuery, lang string) {
	if lang == usr.Language {
		// Because we must answer on every callback request
		_, err := bot.AnswerCallbackQuery(
			tg.NewAnswerCallbackQuery(call.ID),
		)
		errCheck(err)
		return
	}

	err := usr.setLanguage(lang)
	errCheck(err)

	T, err := langSwitch(usr.Language, call.From.LanguageCode)
	errCheck(err)

	go callbackAlert(call, T("message_language_selected"))

	callbackToLanguages(usr, call)
}

func callbackToLanguages(usr *user, call *tg.CallbackQuery) {
	T, err := langSwitch(usr.Language, call.From.LanguageCode)
	errCheck(err)

	text := T("message_language", map[string]interface{}{
		"LanguageCodes": strings.Join(languageTags, "|"),
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = getLanguagesMenuKeyboard(usr)

	_, err = bot.EditMessageText(editText)
	errCheck(err)
}

func getLanguagesMenuKeyboard(usr *user) *tg.InlineKeyboardMarkup {
	T, err := langSwitch(usr.Language)
	errCheck(err)

	var replyMarkup tg.InlineKeyboardMarkup
	for _, tag := range languageTags {
		var this string
		if usr.Language == tag {
			this = switcherStatus
		}

		replyMarkup.InlineKeyboard = append(
			replyMarkup.InlineKeyboard,
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButton(
					fmt.Sprint(languageNames[tag], this),
					fmt.Sprint("switch:language:", tag),
				),
			),
		)
	}
	replyMarkup.InlineKeyboard = append(
		replyMarkup.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)

	return &replyMarkup
}

func callbackUpdateLanguagesKeyboard(usr *user, call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = getLanguagesMenuKeyboard(usr)

	_, err := bot.EditMessageReplyMarkup(&editMarkup)
	errCheck(err)
}
