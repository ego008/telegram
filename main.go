package main

import (
	"encoding/json"
	"flag"
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	appMetrika = make(chan bool)
	bot        *tgbotapi.BotAPI
	config     Configuration
	metrika    botan.Botan
	resNum     = 20 // Select Gelbooru by default, remake in name search(?)
	update     tgbotapi.Update
)

func init() {
	// Read configuration
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("[Configuration] Reading error: %+v", err)
	} else {
		log.Println("[Configuration] Read successfully.")
	}
	// Decode configuration
	if err = json.Unmarshal(file, &config); err != nil {
		log.Fatalf("[Configuration] Decoding error: %+v", err)
	}

	// Initialize bot
	newBot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Panicf("[Bot] Initialize error: %+v", err)
	} else {
		bot = newBot
		log.Printf("[Bot] Authorized as @%s", bot.Self.UserName)
	}

	metrika = botan.New(config.Botan.Token) // Initialize botan
	log.Println("[Botan] ACTIVATED")
}

func main() {
	startUptime := time.Now() // Set start UpTime time

	debugMode := flag.Bool("debug", false, "enable debug logs")
	webhookMode := flag.Bool("webhook", false, "enable webhooks support")
	cacheTime := flag.Int("cache", 0, "cache time in seconds for inline-search results")
	flag.Parse()

	bot.Debug = *debugMode

	updates := make(<-chan tgbotapi.Update)
	if *webhookMode == true {
		if _, err := bot.SetWebhook(tgbotapi.NewWebhook(config.Telegram.Webhook.Set + config.Telegram.Token)); err != nil {
			log.Printf("Set webhook error: %+v", err)
		}
		updates = bot.ListenForWebhook(config.Telegram.Webhook.Listen + config.Telegram.Token)
		go http.ListenAndServe(config.Telegram.Webhook.Serve, nil)
	} else {
		upd := tgbotapi.NewUpdate(0)
		upd.Timeout = 60
		updater, err := bot.GetUpdatesChan(upd)
		if err != nil {
			log.Printf("[Bot] Getting updates error: %+v", err)
		}
		updates = updater
	}

	// Updater
	for update = range updates {
		log.Printf("[Bot] Update response: %+v", updates)

		// Chat actions
		if update.Message != nil {
			switch {
			case update.Message.Command() == "start" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)): // Requirement Telegram platform
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/start", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /start %s", answer.Status)
					appMetrika <- true
				})

				go sendHello(update.Message)

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "help" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)): // Requirement Telegram platform
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/help", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /help %s", answer.Status)
					appMetrika <- true
				})

				go sendHelp(update.Message)

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "cheatsheet" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)):
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/cheatsheet", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /cheatsheet %s", answer.Status)
					appMetrika <- true
				})

				go sendCheatSheet(update.Message)

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "random" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)):
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/random", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /random %s", answer.Status)
					appMetrika <- true
				})

				go sendRandomPost(update.Message)

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "info" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)):
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/info", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /info %s", answer.Status)
					appMetrika <- true
				})

				go sendBotInfo(update.Message, startUptime)

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Chat.IsPrivate() && update.Message.From.ID == config.Telegram.Admin:
				go sendTelegramFileID(update.Message)
			default:
				go getEasterEgg(update.Message) // Secret actions and commands ;)
			}
		}

		// Inline actions
		if update.InlineQuery != nil && len(update.InlineQuery.Query) <= 255 { // Just don't update results if query exceeds the maximum length
			// Track action
			metrika.TrackAsync(update.InlineQuery.From.ID, MetrikaInlineQuery{update.InlineQuery}, "Search", func(answer botan.Answer, err []error) {
				log.Printf("[Botan] Track Search %s", answer.Status)
				appMetrika <- true
			})

			go getInlineResults(update.InlineQuery, *cacheTime)

			<-appMetrika // Send track to Yandex.AppMetrika
		}

		if update.ChosenInlineResult != nil {
			metrika.TrackAsync(update.ChosenInlineResult.From.ID, MetrikaChosenInlineResult{update.ChosenInlineResult}, "Find", func(answer botan.Answer, err []error) {
				log.Printf("[Botan] Track Find %s", answer.Status)
				appMetrika <- true
			})
			<-appMetrika // Send track to Yandex.AppMetrika
		}

		if update.CallbackQuery != nil {
			metrika.TrackAsync(update.ChosenInlineResult.From.ID, MetrikaCallbackQuery{update.CallbackQuery}, "Action", func(answer botan.Answer, err []error) {
				log.Printf("[Botan] Track Find %s", answer.Status)
				appMetrika <- true
			})
			<-appMetrika // Send track to Yandex.AppMetrika
		}
	}
}
