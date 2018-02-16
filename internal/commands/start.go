package commands

import (
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/config"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/models"
	tg "github.com/toby3d/telegram"
)

const queryExample = "hatsune_miku"

func commandStart(msg *tg.Message) {
	usr, err := db.GetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	T, err := i18n.SwitchTo(usr.Language, msg.From.LanguageCode)
	errors.Check(err)

	channelURL := config.Config.UString("telegram.channel.invite")
	if channelURL == "" {
		channelURL = "https://t.me/HentaiDB"
	}

	text := T("message_start", map[string]interface{}{
		"FirstName":         msg.From.FirstName,
		"ChannelURL":        channelURL,
		"Username":          bot.Bot.Self.Username,
		"Query":             queryExample,
		"CommandCheatsheet": models.Cheatsheet,
		"CommandHelp":       models.Help,
	})

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitchSelf(T("button_try"), queryExample),
		),
	)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
