// +build !easterEggs

package main

import (
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func GetEasterEgg(bot *tgbotapi.BotAPI, metrika botan.Botan, update tgbotapi.Update) {
	// Do nothing, because you a not @toby3d
	log.Println("EASTER EGGS: Not today, boy. (¬‿¬)")
}
