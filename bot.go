package main

import (
	"encoding/json"
	"flag"
	b "github.com/botanio/sdk/go"
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/i18n"
	r "gopkg.in/dancannon/gorethink.v2"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	appMetrika = make(chan bool)
	bot        *t.BotAPI
	config     Configuration
	db         *r.Session
	metrika    b.Botan
)

func init() {
	// Read configuration
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("[Configuration] Reading error: %+v", err)
	} else {
		log.Println("[Configuration] Read successfully!")
	}
	// Decode configuration
	if err = json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("[Configuration] Decoding error: %+v", err)
	}

	// Read localization
	i18n.MustLoadTranslationFile("i18n/en-us.all.json")
	i18n.MustLoadTranslationFile("i18n/ru-ru.all.json")
	i18n.MustLoadTranslationFile("i18n/zh-zh.all.json")

	// Initialize RethinkDB
	if rethink, err := r.Connect(r.ConnectOpts{
		Address:  config.DataBase.Address,
		Database: config.DataBase.DataBase,
	}); err != nil {
		log.Fatalln(err)
	} else {
		db = rethink
	}

	// Initialize bot
	teleBot, err := t.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Fatalf("[Bot] Initialize error: %+v", err)
	} else {
		bot = teleBot
		log.Printf("[Bot] Authorized as @%s", bot.Self.UserName)
	}

	metrika = b.New(config.Botan.Token) // Initialize botan
}

func main() {
	debugMode := flag.Bool("debug", false, "enable debug logs")
	webhookMode := flag.Bool("webhook", false, "enable webhooks support")
	cacheTime := flag.Int("cache", 0, "cache time in seconds for inline-search results")
	flag.Parse()

	bot.RemoveWebhook()
	bot.Debug = *debugMode

	updates := make(<-chan t.Update)
	updates = setUpdates(*webhookMode)

	// Updater
	for update := range updates {
		switch {
		case update.Message != nil:
			go getMessages(update.Message)
		case update.InlineQuery != nil && len(update.InlineQuery.Query) <= 255: // Just don't update results if query exceeds the maximum length
			go getInlineResults(*cacheTime, update.InlineQuery)
		case update.ChosenInlineResult != nil:
			go trackInlineResult(update.ChosenInlineResult)
		case update.CallbackQuery != nil:
			go checkCallbackQuery(update.CallbackQuery)
		}
	}
}

func setUpdates(isWebhook bool) <-chan t.Update {
	if isWebhook == true {
		log.Println("[Bot] Webhook activated")
		if _, err := bot.SetWebhook(t.NewWebhook(config.Telegram.Webhook.Set + config.Telegram.Token)); err != nil {
			log.Fatalf("[Bot] Set webhook error: %+v", err)
		}
		go http.ListenAndServe(config.Telegram.Webhook.Serve, nil)
		updates := bot.ListenForWebhook(config.Telegram.Webhook.Listen + config.Telegram.Token)
		return updates
	} else {
		upd := t.NewUpdate(0)
		upd.Timeout = 60
		updates, err := bot.GetUpdatesChan(upd)
		if err != nil {
			log.Fatalf("[Bot] Getting updates error: %+v", err)
		}
		return updates
	}
}
