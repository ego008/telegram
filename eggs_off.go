// +build !eggs

package main

import (
	b "github.com/botanio/sdk/go"
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func init() {
	log.Println("[Easter eggs] Deactivated!")
}

// GetEasterEgg could send an easeter egg. But no.
func easterEggsMessages(message *t.Message) {
	if message.Chat.IsPrivate() || bot.IsMessageToMe(*message) {
		// Track all other messages
		metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "Message", func(answer b.Answer, err []error) {
			log.Printf("[Botan] Track Message %s", answer.Status)
			appMetrika <- true
		})

		<-appMetrika
	}
}

func eggCommand(message *t.Message) {
	if message.Chat.IsPrivate() || bot.IsMessageToMe(*message) {
		// Track all other messages
		metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "Message", func(answer b.Answer, err []error) {
			log.Printf("[Botan] Track Message %s", answer.Status)
			appMetrika <- true
		})

		<-appMetrika
	}
}
