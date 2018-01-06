package main

import tg "github.com/toby3d/telegram"

const queryExample = "hatsune_miku"

func commandStart(msg *tg.Message) {
	usr, err := dbGetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	T, err := langSwitch(usr.Language, msg.From.LanguageCode)
	errCheck(err)

	channelURL := cfg.UString("telegram.channel.invite")
	if channelURL == "" {
		channelURL = "https://t.me/HentaiDB"
	}

	text := T("message_start", map[string]interface{}{
		"FirstName":         msg.From.FirstName,
		"ChannelURL":        channelURL,
		"Username":          bot.Self.Username,
		"Query":             queryExample,
		"CommandCheatsheet": cmdCheatsheet,
		"CommandHelp":       cmdHelp,
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitchSelf(T("button_try"), queryExample),
		),
	)

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
