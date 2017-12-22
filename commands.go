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
	case blackList, whiteList:
		if !msg.HasArgument() {
			return
		}

		listType := blackList
		if listType == whiteList {
			listType = whiteList
		}

		usr, err := dbGetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
		errCheck(err)

		tags := strings.Split(strings.ToLower(msg.CommandArgument()), " ")
		err = usr.addListTags(listType, tags...)
		errCheck(err)

		reply := tg.NewMessage(msg.Chat.ID, "OK")
		_, err = bot.SendMessage(reply)
		errCheck(err)
	}
}
