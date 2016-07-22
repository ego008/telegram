// +build !webhook

package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func init() {
	bot.RemoveWebhook()
	log.Println("Webhook mode deactivated.")
}

func SetUpdater() <-chan tgbotapi.Update {
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	updates, err := bot.GetUpdatesChan(upd)
	if err != nil {
		log.Fatalf("Error getting updates: %+v", err)
	}
	return updates
}
