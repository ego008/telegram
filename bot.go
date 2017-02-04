package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hjson/hjson-go"
	log "github.com/kirillDanshin/dlog"
	f "github.com/valyala/fasthttp"
)

var (
	b       botan.Botan
	bot     *tg.BotAPI
	cfg     map[string]interface{}
	metrika = make(chan bool)
)

func init() {
	// Open configuration file
	config, err := ioutil.ReadFile("config.hjson")
	if err != nil {
		panic(err.Error())
	}

	// Read configuration
	if err = hjson.Unmarshal(config, &cfg); err != nil {
		panic(err.Error())
	}

	b = botan.New(cfg["botan"].(string))

	// Initialize bot
	bot, err = tg.NewBotAPI(cfg["telegram_token"].(string))
	if err != nil {
		panic(err.Error())
	}
	log.F("Authorized as @%s", bot.Self.UserName)
}

func main() {
	debugMode := flag.Bool("debug", false, "enable debug logs")
	webhookMode := flag.Bool("webhook", false, "enable webhooks support")
	cacheTime := flag.Int("cache", 0, "cache time in seconds for inline-search results")
	flag.Parse()

	bot.Debug = *debugMode

	updates := make(<-chan tg.Update)
	updates = setUpdates(*webhookMode)

	// Updater
	for upd := range updates {
		switch {
		case upd.Message != nil:
			go getMessages(upd.Message)
		case upd.InlineQuery != nil && len(upd.InlineQuery.Query) <= 255: // Just don't update results if query exceeds the maximum length
			go getInlineResults(*cacheTime, upd.InlineQuery)
		case upd.ChosenInlineResult != nil:
			go trackInlineResult(upd.ChosenInlineResult)
			// case upd.CallbackQuery != nil:
			// 	go checkCallbackQuery(upd.CallbackQuery)
		}
	}
}

func setUpdates(isWebhook bool) <-chan tg.Update {
	bot.RemoveWebhook()
	if isWebhook {
		if _, err := bot.SetWebhook(
			tg.NewWebhook(
				fmt.Sprintf(cfg["telegram_webhook_set"].(string), cfg["telegram_token"].(string)),
			),
		); err != nil {
			panic(err)
		}
		go f.ListenAndServe(cfg["telegram_webhook_serve"].(string), nil)
		// go http.ListenAndServe(cfg["telegram_webhook_serve"].(string), nil)
		return bot.ListenForWebhook(
			fmt.Sprintf(cfg["telegram_webhook_listen"].(string), cfg["telegram_token"].(string)),
		)
	} else {
		upd := tg.NewUpdate(0)
		upd.Timeout = 60
		updates, err := bot.GetUpdatesChan(upd)
		if err != nil {
			panic(err)
		}
		return updates
	}
}
