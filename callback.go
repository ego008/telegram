package main

/*
import (
	"log"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func checkCallbackQuery(callback *tg.CallbackQuery) {
	switch callback.Data {
	case "nsfw_on":
		locale := checkLanguage(callback.From)
		go switchNSFW(callback.From, true)

		replyMarkup := tg.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(locale("status_on")),
				}), "nsfw_off"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_language"), "to_lang"),
			),
		)

		newKeys := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case "nsfw_off":
		locale := checkLanguage(callback.From)
		go switchNSFW(callback.From, false)

		replyMarkup := tg.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(locale("status_off")),
				}), "nsfw_on"),
			),
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonData(locale("button_language"), "to_lang"),
			),
		)

		newKeys := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
		if _, err := bot.Send(newKeys); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case "to_lang":
		locale := checkLanguage(callback.From)
		replyMarkup := tg.NewInlineKeyboardMarkup(
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
		newKeys := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
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
		newText := tg.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)
		if _, err := bot.Send(newText); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}

func settingsMessage(callback *tg.CallbackQuery) {
	locale := checkLanguage(callback.From)
	nsfw := checkNSFW(callback.From)

	var nsfwBtn tg.InlineKeyboardButton
	if nsfw {
		nsfwBtn = tg.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(locale("status_on")),
		}), "nsfw_off")
	} else {
		nsfwBtn = tg.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(locale("status_off")),
		}), "nsfw_on")
	}

	replyMarkup := tg.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			nsfwBtn,
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData(locale("button_language"), "to_lang"),
		),
	)
	newKeys := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, replyMarkup)
	if _, err := bot.Send(newKeys); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

*/
