package main

import (
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	http "github.com/valyala/fasthttp"
)

var bot *tg.BotAPI

func main() {
	defer db.Close()

	bot, err = tg.NewBotAPI(botToken)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Authorized as @%s", bot.Self.UserName)

	bot.Debug = *flagDebug

	updates := make(<-chan tg.Update)
	updates = setUpdates(*flagWebhook)
	defer bot.RemoveWebhook()

	// Updater
	for upd := range updates {
		switch {
		case upd.Message != nil:
			go message(upd.Message)
		case upd.InlineQuery != nil && len(upd.InlineQuery.Query) <= 255: // Just don't update results if query exceeds the maximum length
			go inline(upd.InlineQuery)
		case upd.ChosenInlineResult != nil:
			go chosenResult(upd.ChosenInlineResult)
		case upd.CallbackQuery != nil:
			go callback(upd.CallbackQuery)
		case upd.ChannelPost != nil:
			go channelPost(upd.ChannelPost)
		}
	}
}

func setUpdates(isWebhook bool) <-chan tg.Update {
	if isWebhook {
		if _, err := bot.SetWebhook(tg.NewWebhook(fmt.Sprintf(webSet, botToken))); err != nil {
			log.Fatalln(err.Error())
		}
		go http.ListenAndServe(webServe, nil)
		return bot.ListenForWebhook(fmt.Sprintf(webListen, botToken))
	}

	upd := tg.NewUpdate(0)
	upd.Timeout = 60
	updates, err := bot.GetUpdatesChan(upd)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return updates
}
