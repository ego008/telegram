// +build webhook

package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
)

func SetUpdater() <-chan tgbotapi.Update {
	log.Println("Webhook mode activated!")
	if _, err := bot.SetWebhook(tgbotapi.NewWebhookWithCert(config.Telegram.WebhookURL+config.Telegram.Token, "cert.pem")); err != nil {
		log.Fatalf("Set webhook error: %+v", err)
	}
	updates := bot.ListenForWebhook(config.Telegram.WebhookPath + config.Telegram.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
	return updates
}
