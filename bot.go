package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	bot        *tgbotapi.BotAPI
	config     Configuration
	locale     Localization
	metrika    botan.Botan
	appMetrika = make(chan bool)
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
	localeFile, err := ioutil.ReadFile("locale.json")
	if err != nil {
		log.Fatalf("[Localization] Reading error: %+v", err)
	} else {
		log.Println("[Localization] Read successfully.")
	}
	// Decode localization
	if err = json.Unmarshal(localeFile, &locale); err != nil {
		log.Fatalf("[Localization] Decoding error: %+v", err)
	}

	// Initialize bot
	teleBot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Fatalf("[Bot] Initialize error: %+v", err)
	} else {
		bot = teleBot
		log.Printf("[Bot] Authorized as @%s", bot.Self.UserName)
	}

	metrika = botan.New(config.Botan.Token) // Initialize botan
}

func main() {
	startUptime := time.Now()

	debugMode := flag.Bool("debug", false, "enable debug logs")
	webhookMode := flag.Bool("webhook", false, "enable webhooks support")
	cacheTime := flag.Int("cache", 0, "cache time in seconds for inline-search results")
	flag.Parse()

	bot.RemoveWebhook()
	bot.Debug = *debugMode

	updates := make(<-chan tgbotapi.Update)
	if *webhookMode == true {
		log.Println("[Bot] Webhook activated")
		if _, err := bot.SetWebhook(tgbotapi.NewWebhook(config.Telegram.Webhook.Set + config.Telegram.Token)); err != nil {
			log.Printf("[Bot] Set webhook error: %+v", err)
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
	for update := range updates {
		switch {
		case update.Message != nil:
			lowerCommand := strings.ToLower(update.Message.Command())
			switch {
			case checkCommand(update.Message, "start"): // Requirement Telegram platform
				messageText := fmt.Sprintf(locale.English.Messages.Start, update.Message.From.FirstName, bot.Self.UserName)
				go sendSimpleMessage(update.Message, lowerCommand, messageText)
			case checkCommand(update.Message, "help"): // Requirement Telegram platform
				go sendSimpleMessage(update.Message, lowerCommand, locale.English.Messages.Help)
			case checkCommand(update.Message, "cheatsheet"):
				go sendSimpleMessage(update.Message, lowerCommand, locale.English.Messages.CheatSheet)
			case checkCommand(update.Message, "random"):
				go sendRandomPost(update.Message)
			case checkCommand(update.Message, "info"):
				go sendBotInfo(update.Message, startUptime)
			case checkCommand(update.Message, "donate"):
				go sendSimpleMessage(update.Message, lowerCommand, locale.English.Messages.Donate)
			case update.Message.Chat.IsPrivate() && update.Message.From.ID == config.Telegram.Admin && update.Message.Text == "":
				go sendTelegramFileID(update.Message) // Admin feature without tracking
			default:
				go getEggMessage(update.Message) // Secret actions and commands ;)
			}
		case update.InlineQuery != nil && len(update.InlineQuery.Query) <= 255: // Just don't update results if query exceeds the maximum length
			go getInlineResults(*cacheTime, update.InlineQuery)
		case update.ChosenInlineResult != nil:
			go sendInlineResult(update.ChosenInlineResult)
			// case update.CallbackQuery != nil:
			// 	go getCallbackAction(update.CallbackQuery)
		}
	}
}

// Any message in private or correct message from groups
func checkCommand(message *tgbotapi.Message, command string) bool {
	isCommand := strings.ToLower(message.Command()) == command
	isPrivate := message.Chat.IsPrivate()
	isGroup := message.Chat.IsGroup()
	isSuperGroup := message.Chat.IsSuperGroup()
	isMessageToMe := bot.IsMessageToMe(*message)
	return isCommand && (isPrivate || ((isGroup || isSuperGroup) && isMessageToMe))
}
