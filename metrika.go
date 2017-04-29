package main

import (
	"log"

	botan "github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	botanCallback struct{ *tg.CallbackQuery }
	botanResult   struct{ *tg.ChosenInlineResult }
	botanMessage  struct{ *tg.Message }
	botanInline   struct{ *tg.InlineQuery }
)

var (
	metrika    botan.Botan
	appMetrika = make(chan bool)
)

func trackCallback(call *tg.CallbackQuery) {
	metrika.TrackAsync(
		call.From.ID,
		botanCallback{call},
		"Callback",
		func(answer botan.Answer, err []error) {
			log.Printf("Track Callback %s", answer.Status)
			appMetrika <- true
		},
	)
}

func trackChosenResult(result *tg.ChosenInlineResult) {
	metrika.TrackAsync(
		result.From.ID,
		botanResult{result},
		"Find",
		func(answer botan.Answer, err []error) {
			log.Printf("Track Find %s", answer.Status)
			appMetrika <- true
		},
	)
}

func trackMessage(msg *tg.Message, label string) {
	metrika.TrackAsync(
		msg.From.ID,
		botanMessage{msg},
		label,
		func(answer botan.Answer, err []error) {
			log.Printf("Track %s %s", label, answer.Status)
			appMetrika <- true
		},
	)
}

func trackInline(inline *tg.InlineQuery) {
	metrika.TrackAsync(
		inline.From.ID,
		botanInline{inline},
		"Search",
		func(answer botan.Answer, err []error) {
			log.Printf("Track Search %s", answer.Status)
			appMetrika <- true
		},
	)
}
