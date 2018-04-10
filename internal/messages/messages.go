package messages

import (
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/commands"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
	"strings"
)

func Message(msg *tg.Message) {
	// Got message from myself
	if bot.Bot.IsMessageFromMe(msg) {
		return
	}

	if bot.Bot.IsCommandToMe(msg) {
		log.Ln("IsCommandToMe", true)
		commands.Commands(msg)
		return
	}

	if msg.IsText() && bot.Bot.IsReplyToMe(msg) {
		log.D(msg.ReplyToMessage)
		usr, err := database.GetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
		errors.Check(err)

		T, err := i18n.SwitchTo(user.Locale, msg.From.LanguageCode)
		errors.Check(err)

		var listType string
		switch {
		case strings.EqualFold(
			msg.ReplyToMessage.Text,
			T("message_input_blacklist_tags", map[string]interface{}{"Limit": 5}),
		):
			listType = models.BlackList
		case strings.EqualFold(
			msg.ReplyToMessage.Text,
			T("message_input_whitelist_tags", map[string]interface{}{"Limit": 5}),
		):
			listType = models.WhiteList
		default:
			return
		}

		tags := strings.Split(strings.ToLower(msg.Text), " ")
		err = database.AddListTags(usr, listType, tags...)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, "OK")
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
	}

	log.D(*msg)
}
