package main

import tg "github.com/toby3d/go-telegram"

func commandHelp(msg *tg.Message) {
	_, err := bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	T, err := langSwitch(msg.From.LanguageCode)
	errCheck(err)

	text := T("message_help", map[string]interface{}{
		"CommandSettings":   cmdSettings,
		"CommandCheatsheet": cmdCheatsheet,
		"CommandPatreon":    cmdPatreon,
		"CommandInfo":       cmdInfo,
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
