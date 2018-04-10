package commands

import (
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func commandHelp(msg *tg.Message) {
	usr, err := database.GetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale, msg.From.LanguageCode)
	errors.Check(err)

	text := T("message_help", map[string]interface{}{
		"CommandSettings":   models.Settings,
		"CommandCheatsheet": models.Cheatsheet,
		"CommandPatreon":    models.Patreon,
		"CommandInfo":       models.Info,
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
