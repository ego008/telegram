package main

import (
	"fmt"

	"github.com/HentaiDB/HentaiDBot/internal/config"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

var channel = make(chan tg.Update, 100)

func getUpdatesChannel() tg.UpdatesChannel {
	log.Ln("getUpdatesChannel")
	if !*flagWebhook {
		log.Ln("Remove old webhook...")
		_, err := bot.DeleteWebhook()
		errors.Check(err)

		log.Ln("Create LongPolling updates channel...")
		return bot.NewLongPollingChannel(&tg.GetUpdatesParameters{
			Offset:  0,
			Limit:   100,
			Timeout: 60,
		})
	}

	set := config.Config.UString("telegram.webhook.set")
	listen := config.Config.UString("telegram.webhook.listen")
	serve := config.Config.UString("telegram.webhook.serve")

	log.Ln("Trying set webhook on", fmt.Sprint(set, listen, bot.AccessToken))

	log.Ln("Create new webhook...")
	webhook := tg.NewWebhook(fmt.Sprint(set, listen, bot.AccessToken), nil)
	webhook.MaxConnections = 40

	log.Ln("Create Webhook updates channel...")
	return bot.NewWebhookChannel(
		webhook,
		"", "",
		set, fmt.Sprint(listen, bot.AccessToken), serve,
	)
}
