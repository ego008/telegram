// +build !eggs

package main

import (
	"log"

	"github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	log.Println("[Easter eggs] Deactivated!")
}

// GetEasterEgg could send an easeter egg. But no.
func easterEggsMessages(message *tg.Message) {
	if message.Chat.IsPrivate() || bot.IsMessageToMe(*message) {
		// Track all other messages
		b.TrackAsync(message.From.ID, MetrikaMessage{message}, "Message", func(answer botan.Answer, err []error) {
			log.Printf("[Botan] Track Message %s", answer.Status)
			metrika <- true
		})

		<-metrika
	}
}

func eggCommand(message *tg.Message) {
	if message.Chat.IsPrivate() || bot.IsMessageToMe(*message) {
		// Track all other messages
		b.TrackAsync(message.From.ID, MetrikaMessage{message}, "Message", func(answer botan.Answer, err []error) {
			log.Printf("[Botan] Track Message %s", answer.Status)
			metrika <- true
		})

		<-metrika
	}
}
