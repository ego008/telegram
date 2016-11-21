package main

import (
	"fmt"
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

func checkCallbackQuery(callback *t.CallbackQuery) {
	log.Printf("Callback: %#v", callback)
	switch callback.Data {
	case "nsfw_on":
		switchNSFW(callback.From, true)

		replyMarkup := t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData("🔞 NSFW ON", "nsfw_off"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale.English.Buttons.Language, "to_lang"),
			),
		)

		newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}

	case "nsfw_off":
		switchNSFW(callback.From, false)

		replyMarkup := t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData("🔞 NSFW OFF", "nsfw_on"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale.English.Buttons.Language, "to_lang"),
			),
		)

		newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}

	case "to_lang":
		replyMarkup := t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale.English.Name, "lang_english"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale.Russian.Name, "lang_russian"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale.TChinese.Name, "lang_tchinese"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale.SChinese.Name, "lang_schinese"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale.English.Buttons.Cancel, "to_settings"),
			),
		)
		newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case "to_settings":
		settingsMessage(callback)
	}

	if strings.HasPrefix(callback.Data, "lang_") {
		lang := strings.TrimLeft(callback.Data, "lang_")
		changeLanguage(callback.From, lang)
		settingsMessage(callback)
	}
}

func settingsMessage(callback *t.CallbackQuery) {
	lang := checkLanguage(callback.From)
	nsfw := checkNSFW(callback.From)

	var nsfwBtn t.InlineKeyboardButton
	if nsfw {
		nsfwBtn = t.NewInlineKeyboardButtonData("🔞 NSFW ON", "nsfw_off")
	} else {
		nsfwBtn = t.NewInlineKeyboardButtonData("🔞 NSFW OFF", "nsfw_on")
	}

	replyMarkup := t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(nsfwBtn),
		t.NewInlineKeyboardRow(t.NewInlineKeyboardButtonData(locale.English.Buttons.Language, "to_lang")),
	)
	newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)

	text := fmt.Sprintf("%s: %s", locale.English.Buttons.Language, lang)
	newText := t.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)

	if _, err := bot.Send(newText); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
	if _, err := bot.Send(newKeys); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}
