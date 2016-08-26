// +build !easterEggs

package main

import (
	"github.com/botanio/sdk/go"
	"log"
)

func init() {
	log.Println("[Easter eggs] DEACTIVATED")
}

// GetEasterEgg could send an easeter egg. But no.
func getEasterEgg() {
	switch {
	case update.Message.Chat.IsPrivate() || (update.Message.Chat.ID == config.Telegram.SuperGroup && bot.IsMessageToMe(*update.Message)):
		// Track all other messages
		metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "Message", func(answer botan.Answer, err []error) {
			log.Printf("[Botan] Track Message %s", answer.Status)
			appMetrika <- true
		})
		<-appMetrika
	default:
		// If Message from ofiicial group - skip trash tracking data
		log.Println("[Botan] Skip Message in official group")
	}
}
