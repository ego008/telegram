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
func GetEasterEgg() {
	// Track all other messages
	metrika.TrackAsync(update.Message.From.ID, update.Message, "Message", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Message: %+v", answer)
		appMetrika <- true
	})
	<-appMetrika
}
