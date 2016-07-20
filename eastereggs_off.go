// +build !easterEggs

package main

import (
	"log"

	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// GetEasterEgg could send an easeter egg. But no.
func GetEasterEgg(bot *tgbotapi.BotAPI, metrika botan.Botan, update tgbotapi.Update) {
	// Do nothing, because you're not @toby3d
	log.Println("EASTER EGGS: Not today, boy. (¬‿¬)")
}
