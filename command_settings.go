package main

import tg "github.com/toby3d/go-telegram"

func commandSettings(msg *tg.Message) {
	_, err := bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	T, err := langSwitch(msg.From.LanguageCode)
	errCheck(err)

	text := T("message_settings", map[string]interface{}{
		"Language":  langTags[msg.From.LanguageCode],
		"Resources": "Gelbooru",
		"Ratings":   "",
		"Blacklist": "",
		"Whitelist": "",
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_language"), "menu:settings:language"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_resources"), "menu:settings:resources"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_ratings"), "menu:settings:ratings"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_blacklist"), "menu:settings:blacklist"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_whitelist"), "menu:settings:whitelist"),
		),
	)

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
