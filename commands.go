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

func commands(msg *tg.Message) {
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
	case blackList:
		if !msg.HasArgument() {
			return
		}

		usr, err := dbGetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
		errCheck(err)

		tags := strings.Split(strings.ToLower(msg.CommandArgument()), " ")
		for _, tag := range tags {
			for i := range usr.Blacklist {
				if usr.Blacklist[i] == tag {
					tags = append(tags[:i], tags[i+1:]...)
					continue
				}

			}
		}

		err = usr.addListTags(blackList, tags[1:]...)
		errCheck(err)
	}
}
