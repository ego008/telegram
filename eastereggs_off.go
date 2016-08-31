// +build !easterEggs

package main

import (
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func init() {
	log.Println("[Easter eggs] DEACTIVATED")
}

// GetEasterEgg could send an easeter egg. But no.
func getEasterEgg(message *tgbotapi.Message) {
	switch {
	case message.Chat.IsPrivate() || (message.Chat.ID == config.Telegram.SuperGroup && bot.IsMessageToMe(*message)):
		// Track all other messages
		metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "Message", func(answer botan.Answer, err []error) {
			log.Printf("[Botan] Track Message %s", answer.Status)
			appMetrika <- true
		})
		<-appMetrika
	default:
		// If Message from ofiicial group - skip trash tracking data
		log.Println("[Botan] Skip Message in official group")
	}
}
