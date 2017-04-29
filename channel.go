package main

import (
	"log"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func channelPost(msg *tg.Message) {
	if msg.Chat.ID == int64(chID) {
		ifDev := strings.Contains(msg.Text, "#dev")
		ifNews := strings.Contains(msg.Text, "#news")
		ifBot := strings.Contains(msg.Text, "#bot")
		if msg.Text != "" && ifDev && ifNews && ifBot {
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
