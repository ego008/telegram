// +build !easterEggs

package main

import (
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func init() {
	log.Println("[Easter eggs] Deactivated")
}

// GetEasterEgg could send an easeter egg. But no.
func getEggMessage(message *tgbotapi.Message, metrika botan.Botan, appMetrika chan bool) {
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
	}
}
