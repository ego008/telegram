// +build webhook

package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func init() {
	bot.RemoveWebhook()
	log.Println("Webhook mode activated!")
}

func SetUpdater() <-chan tgbotapi.Update {
	if _, err := bot.SetWebhook(tgbotapi.NewWebhook(config.Telegram.Webhook.Set + config.Telegram.Token)); err != nil {
		log.Fatalf("Set webhook error: %+v", err)
	}
	updates := bot.ListenForWebhook(config.Telegram.Webhook.Listen + config.Telegram.Token)
	return updates
}
