package commands

import (
	"fmt"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/models"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func Commands(msg *tg.Message) {
	cmd := strings.ToLower(msg.Command())
	log.Ln("/" + cmd)

	cmd = strings.TrimSuffix(cmd, fmt.Sprint("@", strings.ToLower(bot.Bot.Self.Username)))

	switch cmd {
	case models.Start:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandStart(msg)
	case models.Help:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandHelp(msg)
	case models.Settings:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandSettings(msg)
	case models.Cheatsheet:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandCheatsheet(msg)
	case models.Random:
		commandRandom(msg)
	case models.BlackList,
		models.WhiteList:
		if !msg.HasCommandArgument() ||
			!msg.Chat.IsPrivate() {
			return
		}

		usr, err := db.GetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
		errors.Check(err)

		tags := strings.Split(strings.ToLower(msg.CommandArgument()), " ")
		err = db.AddListTags(usr, cmd, tags...)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, "OK")
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
	}
}
