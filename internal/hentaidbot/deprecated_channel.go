package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func channelPost(msg *tg.Message) {
	if msg.Chat.ID == int64(chID) {
		ifNews := false
		ifBot := false
		ifUserName := false

		if msg.Text == "" {
			if msg.Caption == "" {
				return
			}
			ifNews = strings.Contains(msg.Caption, "#news")
			ifBot = strings.Contains(msg.Caption, "#bot")
			ifUserName = strings.Contains(msg.Caption, fmt.Sprint("@", bot.Self.UserName))
		} else {
			ifNews = strings.Contains(msg.Text, "#news")
			ifBot = strings.Contains(msg.Text, "#bot")
			ifUserName = strings.Contains(msg.Text, fmt.Sprint("@", bot.Self.UserName))
		}

		if ifNews && (ifBot || ifUserName) {
			users, err := getUsers()
			if err != nil {
				log.Println(err.Error())
				return
			}

			for _, user := range users {
				forw := tg.NewForward(int64(user.ID), int64(chID), msg.MessageID)
				if _, err := bot.Send(forw); err != nil {
					log.Println(err.Error())
				}
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}
