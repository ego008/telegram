// +build !eggs

package main

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	log.Println("Easter Eggs deactivated!")
}

func msgEasterEgg(msg *tg.Message) {
	trackMessage(msg, "Message")
	<-appMetrika
}

func cmdEasterEgg(msg *tg.Message) {
	trackMessage(msg, "Message")
	<-appMetrika
}
