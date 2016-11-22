package main

import (
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

func checkCallbackQuery(callback *t.CallbackQuery) {
	switch callback.Data {
	case "nsfw_on":
		locale := checkLanguage(callback.From)
		go switchNSFW(callback.From, true)

		replyMarkup := t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(locale("status_on")),
				}), "nsfw_off"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_language"), "to_lang"),
			),
		)

		newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case "nsfw_off":
		locale := checkLanguage(callback.From)
		go switchNSFW(callback.From, false)

		replyMarkup := t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(locale("status_off")),
				}), "nsfw_on"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_language"), "to_lang"),
			),
		)

		newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case "to_lang":
		locale := checkLanguage(callback.From)
		replyMarkup := t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData("🇬🇧 English", "lang_en-us"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru-ru"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData("🇹🇼 正體中文", "lang_zh-zh"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_cancel"), "to_settings"),
			),
		)
		newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case "to_settings":
		go settingsMessage(callback)
	}

	if strings.HasPrefix(callback.Data, "lang_") {
		lang := strings.TrimLeft(callback.Data, "lang_")
		locale := changeLanguage(callback.From, lang)

		go settingsMessage(callback)

		text := locale("message_settings")
		newText := t.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)
		if _, err := bot.Send(newText); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}

func settingsMessage(callback *t.CallbackQuery) {
	locale := checkLanguage(callback.From)
	nsfw := checkNSFW(callback.From)

	var nsfwBtn t.InlineKeyboardButton
	if nsfw {
		nsfwBtn = t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(locale("status_on")),
		}), "nsfw_off")
	} else {
		nsfwBtn = t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(locale("status_off")),
		}), "nsfw_on")
	}

	replyMarkup := t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			nsfwBtn,
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData(locale("button_language"), "to_lang"),
		),
	)
	newKeys := t.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
	if _, err := bot.Send(newKeys); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}
