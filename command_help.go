package main

import tg "github.com/toby3d/telegram"

func commandHelp(msg *tg.Message) {
	usr, err := dbGetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	T, err := langSwitch(usr.Language, msg.From.LanguageCode)
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
