package main

import (
	"log"

	tg "github.com/toby3d/telegram"
)

var bot *tg.Bot

func main() {
	defer func() {
		err := db.Close()
		errCheck(err)
	}()

	var err error
	bot, err = tg.NewBot(cfg.UString("telegram.token"))
	if err != nil {
		log.Fatalln(err.Error())
	}
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
			callbackQuery(update.CallbackQuery)
		case update.ChannelPost != nil:
			// channelPost(update.ChannelPost)
		default:
			continue
		}
	}
}
