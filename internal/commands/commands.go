package commands

import (
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func Commands(msg *tg.Message) {
	switch {
	case msg.IsCommand(models.CommandStart):
		if !msg.Chat.IsPrivate() {
			return
		}

		commandStart(msg)
	case msg.IsCommand(models.CommandHelp):
		if !msg.Chat.IsPrivate() {
			return
		}

		commandHelp(msg)
	case msg.IsCommand(models.CommandSettings):
		if !msg.Chat.IsPrivate() {
			return
		}

		commandSettings(msg)
	case msg.IsCommand(models.CommandCheatsheet):
		if !msg.Chat.IsPrivate() {
			return
		}

		commandCheatsheet(msg)
	case msg.IsCommand(models.CommandRandom):
		commandRandom(msg)
	case msg.IsCommand(models.CommandBlackList),
		models.WhiteList:
		if !msg.HasCommandArgument() ||
			!msg.Chat.IsPrivate() {
			return
		}

		usr, err := database.GetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
		errors.Check(err)

		tags := strings.Split(strings.ToLower(msg.CommandArgument()), " ")
		err = database.AddListTags(usr, cmd, tags...)
		errors.Check(err)

		reply := tg.NewMessage(msg.Chat.ID, "OK")
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
	}
}
