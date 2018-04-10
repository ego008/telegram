package commands

import (
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/callbacks"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	tg "github.com/toby3d/telegram"
)

func commandSettings(msg *tg.Message) {
	usr, err := database.GetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale, msg.From.LanguageCode)
	errors.Check(err)

	var activeRes []string
	for k, v := range usr.Resources {
		if !v {
			continue
		}

		title := resources.Resources[k].UString("title")
		activeRes = append(activeRes, title)
	}

	ratings, err := usr.GetRatingsStatus()
	errors.Check(err)

	text := T("message_settings", map[string]interface{}{
		"Language":  i18n.Names[user.Locale],
		"Resources": strings.Join(activeRes, ", "),
		"Ratings":   ratings,
		"Blacklist": strings.Join(user.Blacklist, "`, `"),
		"Whitelist": strings.Join(usr.Whitelist, "`, `"),
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = callbacks.GetSettingsMenuKeyboard(usr)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
