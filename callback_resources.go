package main

import (
	"fmt"

	// log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func callbackToggleResource(usr *user, call *tg.CallbackQuery, res string) {
	err := usr.toggleResource(res)
	errCheck(err)

	callbackUpdateResourcesKeyboard(usr, call)
}

func callbackToResources(usr *user, call *tg.CallbackQuery) {
	T, err := langSwitch(usr.Language, call.From.LanguageCode)
	errCheck(err)

	text := T("message_resources")

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = getResourcesMenuKeyboard(usr)

	_, err = bot.EditMessageText(editText)
	errCheck(err)
}

func getResourcesMenuKeyboard(usr *user) *tg.InlineKeyboardMarkup {
	T, err := langSwitch(usr.Language)
	errCheck(err)

	var row int
	var replyMarkup tg.InlineKeyboardMarkup
	for i, tag := range resourcesTags {
		if i%2 == 0 {
			replyMarkup.InlineKeyboard = append(
				replyMarkup.InlineKeyboard, tg.NewInlineKeyboardRow(),
			)
			row++
		}

		replyMarkup.InlineKeyboard[row-1] = append(
			replyMarkup.InlineKeyboard[row-1],
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Resources[tag]],
					resources[tag].UString("title"),
				),
				fmt.Sprint("toggle:resource:", tag),
			),
		)

		i++
	}
	replyMarkup.InlineKeyboard = append(
		replyMarkup.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)

	return &replyMarkup
}

func callbackUpdateResourcesKeyboard(usr *user, call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = getResourcesMenuKeyboard(usr)

	_, err := bot.EditMessageReplyMarkup(&editMarkup)
	errCheck(err)
}
