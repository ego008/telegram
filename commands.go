package main

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

const (
	cmdStart      = "start"
	cmdHelp       = "help"
	cmdSettings   = "settings"
	cmdCheatsheet = "cheatsheet"
	cmdInfo       = "info"
	cmdPatreon    = "patreon"
	cmdRandom     = "random"
)

func commands(msg *tg.Message) {
	cmd := strings.ToLower(msg.Command())
	log.Ln("/" + cmd)

	cmd = strings.TrimSuffix(cmd, fmt.Sprint("@", strings.ToLower(bot.Self.Username)))

	switch cmd {
	case cmdStart:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandStart(msg)
	case cmdHelp:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandHelp(msg)
	case cmdSettings:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandSettings(msg)
	case cmdCheatsheet:
		if !msg.Chat.IsPrivate() {
			return
		}

		commandCheatsheet(msg)
	case cmdRandom:
		commandRandom(msg)
	case blackList, whiteList:
		if !msg.HasArgument() ||
			!msg.Chat.IsPrivate() {
			return
		}

		usr, err := dbGetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
		errCheck(err)

		tags := strings.Split(strings.ToLower(msg.CommandArgument()), " ")
		err = usr.addListTags(cmd, tags...)
		errCheck(err)

		reply := tg.NewMessage(msg.Chat.ID, "OK")
		_, err = bot.SendMessage(reply)
		errCheck(err)
	}
}
