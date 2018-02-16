package main

import (
	"github.com/HentaiDB/HentaiDBot/internal/commands"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func message(msg *tg.Message) {
	// Getted message from myself
	if msg.From.ID == bot.Self.ID {
		return
	}

	// Getted forwarded message from myself
	if msg.ForwardFrom != nil {
		if msg.ForwardFrom.ID == bot.Self.ID {
			return
		}
	}

	if !msg.IsCommand("") {
		log.D(*msg)
		return
	}

	commands.Commands(msg)
}
