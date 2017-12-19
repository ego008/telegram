package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/go-telegram"
)

const (
	cmdStart      = "start"
	cmdHelp       = "help"
	cmdSettings   = "settings"
	cmdCheatsheet = "cheatsheet"
	cmdInfo       = "info"
	cmdPatreon    = "patreon"
)

func command(msg *tg.Message) {
	cmd := strings.ToLower(msg.Command())
	log.Ln("/" + cmd)

	switch cmd {
	case cmdStart:
		commandStart(msg)
	case cmdHelp:
		commandHelp(msg)
	case cmdSettings:
		commandSettings(msg)
	case cmdCheatsheet:
		commandCheatsheet(msg)
	}
}
