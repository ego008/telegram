package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/botanio/sdk/go"
	"github.com/hjson/hjson-go"
	log "github.com/kirillDanshin/dlog"
	f "github.com/valyala/fasthttp"
	tg "gopkg.in/telegram-bot-api.v4"
)

var (
	debugMode   = flag.Bool("debug", false, "enable debug logs")
	webhookMode = flag.Bool("webhook", false, "enable webhooks support")
	cacheTime   = flag.Int("cache", 0, "cache time in seconds for inline-search results")

	b       botan.Botan
	bot     *tg.BotAPI
	cfg     map[string]interface{}
	metrika = make(chan bool)
)

func init() {
	flag.Parse() // Parse flags

	// Open configuration file
	config, err := ioutil.ReadFile("config.hjson")
	if err != nil {
		panic(err.Error())
	}
	if err = hjson.Unmarshal(config, &cfg); err != nil {
		panic(err.Error())
	}

	log.Ln("TRY RUNNING VERSION", cfg["telegram_version_name"].(string))

	b = botan.New(cfg["botan"].(string)) // Set metrika counter

	// Initialize bot
	bot, err = tg.NewBotAPI(cfg["telegram_token"].(string))
	if err != nil {
		panic(err.Error())
	}
	log.F("Authorized as @%s", bot.Self.UserName)
}

func main() {
	bot.Debug = *debugMode

	updates := make(<-chan tg.Update)
	updates = SetUpdates(*webhookMode)
	defer bot.RemoveWebhook()

	// Updater
	for upd := range updates {
		switch {
		case upd.Message != nil:
			go GetMessage(upd.Message)
		case upd.InlineQuery != nil && len(upd.InlineQuery.Query) <= 255: // Just don't update results if query exceeds the maximum length
			go GetInlineResults(upd.InlineQuery)
		case upd.ChosenInlineResult != nil:
			go TrackChosenInlineResult(upd.ChosenInlineResult)
		case upd.CallbackQuery != nil:
			go CheckCallbackQuery(upd.CallbackQuery)
		}
	}
}

func SetUpdates(isWebhook bool) <-chan tg.Update {
	if isWebhook {
		if _, err := bot.SetWebhook(
			tg.NewWebhook(
				fmt.Sprintf(cfg["telegram_webhook_set"].(string), cfg["telegram_token"].(string)),
			),
		); err != nil {
			panic(err)
		}
		go f.ListenAndServe(cfg["telegram_webhook_serve"].(string), nil)
		return bot.ListenForWebhook(
			fmt.Sprintf(cfg["telegram_webhook_listen"].(string), cfg["telegram_token"].(string)),
		)
	}

	upd := tg.NewUpdate(0)
	upd.Timeout = 60
	updates, err := bot.GetUpdatesChan(upd)
	if err != nil {
		panic(err)
	}
	return updates
}
