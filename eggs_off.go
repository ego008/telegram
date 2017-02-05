// +build !eggs

package main

import (
	"github.com/botanio/sdk/go"
	log "github.com/kirillDanshin/dlog"
	tg "gopkg.in/telegram-bot-api.v4"
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
