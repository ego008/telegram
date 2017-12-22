package main

import (
	"fmt"

	// log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/go-telegram"
)

func callbackToList(usr *user, call *tg.CallbackQuery, listType string) {
	T, err := langSwitch(usr.Language, call.From.LanguageCode)
	errCheck(err)

	text := T(fmt.Sprint("message_", listType), map[string]interface{}{
		"CommandCheatsheet": cmdCheatsheet,
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = getListMenuKeyboard(usr, listType)

	_, err = bot.EditMessageText(editText)
	errCheck(err)
}

func getListMenuKeyboard(usr *user, listType string) *tg.InlineKeyboardMarkup {
	T, err := langSwitch(usr.Language)
	errCheck(err)

	var tags []string
	switch listType {
	case blackList:
		tags = usr.Blacklist
	case whiteList:
		tags = usr.Whitelist
	}

	replyMarkup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_tags_add"), fmt.Sprint("add:tags:", listType),
			),
		),
	)
	row := 1
	for i, tag := range tags {
		if i%2 == 0 {
			replyMarkup.InlineKeyboard = append(
				replyMarkup.InlineKeyboard, tg.NewInlineKeyboardRow(),
			)
			row++
		}

		replyMarkup.InlineKeyboard[row-1] = append(
			replyMarkup.InlineKeyboard[row-1],
			tg.NewInlineKeyboardButton(
				tag, fmt.Sprint("remove:", listType, ":", tag),
			),
		)
	}
	replyMarkup.InlineKeyboard = append(
		replyMarkup.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)

	return replyMarkup
}

func callbackUpdateListKeyboard(usr *user, call *tg.CallbackQuery, listType string) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = getListMenuKeyboard(usr, listType)

	_, err := bot.EditMessageReplyMarkup(&editMarkup)
	errCheck(err)
}
