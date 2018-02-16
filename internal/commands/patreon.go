package commands

import (
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	tg "github.com/toby3d/telegram"
)

func commandPatreon(msg *tg.Message) {
	usr, err := db.GetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	T, err := i18n.SwitchTo(usr.Language, msg.From.LanguageCode)
	errors.Check(err)

	text := T("message_patreon", 0, map[string]interface{}{
		"Patrons": "",
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
