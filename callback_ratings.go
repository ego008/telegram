package main

import (
	"fmt"
	"strings"

	// log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/go-telegram"
)

func callbackToggleRating(usr *user, call *tg.CallbackQuery, rating string) {
	var err error
	switch rating {
	case "safe":
		err = usr.toggleSafe()
	case "questionable":
		err = usr.toggleQuestionable()
	case "explicit":
		err = usr.toggleExplicit()
	}
	errCheck(err)

	callbackUpdateRatingsKeyboard(usr, call)
}

func callbackToRatings(usr *user, call *tg.CallbackQuery) {
	T, err := langSwitch(usr.Language, call.From.LanguageCode)
	errCheck(err)

	text := T("message_ratings", map[string]interface{}{
		"Safe":         strings.Title(T("rating_safe")),
		"Questionable": strings.Title(T("rating_questionable")),
		"Explicit":     strings.Title(T("rating_explicit")),
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = getRatingsMenuKeyboard(usr)

	_, err = bot.EditMessageText(editText)
	errCheck(err)
}

func getRatingsMenuKeyboard(usr *user) *tg.InlineKeyboardMarkup {
	T, err := langSwitch(usr.Language)
	errCheck(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprintln(
					toggleStatus[usr.Ratings.Safe],
					strings.Title(T("rating_safe")),
				),
				"toggle:rating:safe",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprintln(
					toggleStatus[usr.Ratings.Questionable],
					strings.Title(T("rating_questionable")),
				),
				"toggle:rating:questionable",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprintln(
					toggleStatus[usr.Ratings.Exlplicit],
					strings.Title(T("rating_explicit")),
				),
				"toggle:rating:explicit",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)
}

func callbackUpdateRatingsKeyboard(usr *user, call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = getRatingsMenuKeyboard(usr)

	_, err := bot.EditMessageReplyMarkup(&editMarkup)
	errCheck(err)
}
