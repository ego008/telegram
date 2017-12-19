package main

import tg "github.com/toby3d/go-telegram"

func commandPatreon(msg *tg.Message) {
	_, err := bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	T, err := langSwitch(msg.From.LanguageCode)
	errCheck(err)

	text := T("message_patreon", 0, map[string]interface{}{
		"Patrons": "",
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
