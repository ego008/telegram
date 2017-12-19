package main

import tg "github.com/toby3d/go-telegram"

func commandCheatsheet(msg *tg.Message) {
	_, err := bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	T, err := langSwitch(msg.From.LanguageCode)
	errCheck(err)

	text := T("message_cheatsheet")

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.DisableWebPagePreview = true

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
