package main

import (
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func messages(msg *tg.Message) {
	// Getted message from myself
	if bot.IsMessageFromMe(msg) {
		return
	}

	log.Ln("IsMessageFromMe", false)

	if bot.IsCommandToMe(msg) {
		log.Ln("IsCommandToMe", true)
		commands(msg)
		return
	}

	log.Ln("IsCommandToMe", false)

	//log.D(*msg)
}
