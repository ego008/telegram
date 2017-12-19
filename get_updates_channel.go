package main

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/go-telegram"
)

var channel = make(chan tg.Update, 100)

func getUpdatesChannel() tg.UpdatesChannel {
	log.Ln("getUpdatesChannel")
	if !*flagWebhook {
		log.Ln("Remove old webhook...")
		_, err := bot.DeleteWebhook()
		errCheck(err)

		log.Ln("Create LongPolling updates channel...")
		return bot.NewLongPollingChannel(&tg.GetUpdatesParameters{
			Offset:  0,
			Limit:   100,
			Timeout: 60,
		})
	}

	set := cfg.UString("telegram.webhook.set")
	listen := cfg.UString("telegram.webhook.listen")
	serve := cfg.UString("telegram.webhook.serve")

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
