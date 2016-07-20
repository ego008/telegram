// +build !webhook

package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func SetUpdater() <-chan tgbotapi.Update {
	log.Println("Webhook mode deactivated.")
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	updates, err := bot.GetUpdatesChan(upd)
	if err != nil {
		log.Fatalf("Error getting updates: %+v", err)
	}
	return updates
}
