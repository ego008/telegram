package main

import (
	"fmt"
	"strings"

	// log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func callbackToggleTypes(usr *user, call *tg.CallbackQuery, resultType string) {
	var err error
	switch resultType {
	case "image":
		err = usr.toggleTypeImage()
	case "animation":
		err = usr.toggleTypeAnimation()
	case "video":
		err = usr.toggleTypeVideo()
	}
	errCheck(err)

	callbackUpdateTypesKeyboard(usr, call)
}

func callbackToTypes(usr *user, call *tg.CallbackQuery) {
	// T, err := langSwitch(usr.Language, call.From.LanguageCode)
	// errCheck(err)

	text := "Here you can select types of content which you want see in results."

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = getTypesMenuKeyboard(usr)

	_, err := bot.EditMessageText(editText)
	errCheck(err)
}

func getTypesMenuKeyboard(usr *user) *tg.InlineKeyboardMarkup {
	T, err := langSwitch(usr.Language)
	errCheck(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Types.Image],
					strings.Title(T("type_image")),
				),
				"toggle:type:image",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Types.Animation],
					strings.Title(T("type_animation")),
				),
				"toggle:type:animation",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Types.Video],
					strings.Title(T("type_video")),
				),
				"toggle:type:video",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)
}

func callbackUpdateTypesKeyboard(usr *user, call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = getTypesMenuKeyboard(usr)

	_, err := bot.EditMessageReplyMarkup(&editMarkup)
	errCheck(err)
}
