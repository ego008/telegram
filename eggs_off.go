// +build !eggs

package main

import (
	"github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/kirillDanshin/dlog"
)

func init() {
	log.Ln("[Easter eggs] Deactivated!")
}

// GetEasterEgg could send an easeter egg. But no.
func EasterEggsMessages(message *tg.Message) {
	if message.Chat.IsPrivate() || bot.IsMessageToMe(*message) {
		// Track all other messages
		b.TrackAsync(message.From.ID, struct{ *tg.Message }{message}, "Message", func(answer botan.Answer, err []error) {
			log.Ln("Track Message", answer.Status)
			metrika <- true
		})

		<-metrika
	}
}

func EggCommand(message *tg.Message) {
	if message.Chat.IsPrivate() || bot.IsMessageToMe(*message) {
		// Track all other messages
		b.TrackAsync(message.From.ID, struct{ *tg.Message }{message}, "Message", func(answer botan.Answer, err []error) {
			log.Ln("Track Message", answer.Status)
			metrika <- true
		})

		<-metrika
	}
}
