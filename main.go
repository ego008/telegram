package main

import (
	"log"

	"github.com/HentaiDB/HentaiDBot/internal/callbacks"
	"github.com/HentaiDB/HentaiDBot/internal/config"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	tg "github.com/toby3d/telegram"
)

var bot *tg.Bot

func main() {
	defer func() {
		err := db.DB.Close()
		errors.Check(err)
	}()

	var err error
	bot, err = tg.NewBot(config.Config.UString("telegram.token"))
	errors.Check(err)
	log.Print("Authorized as @", bot.Self.Username)

	// Updater
	for update := range getUpdatesChannel() {
		switch {
		case update.Message != nil:
			messages(update.Message)
		case update.InlineQuery != nil &&
			len(update.InlineQuery.Query) <= 255:
			inlineQuery(update.InlineQuery)
		case update.ChosenInlineResult != nil:
			// ChosenInlineResult(update.ChosenInlineResult)
		case update.CallbackQuery != nil:
			callbacks.CallbackQuery(update.CallbackQuery)
		case update.ChannelPost != nil:
			// channelPost(update.ChannelPost)
		}
	}
}
