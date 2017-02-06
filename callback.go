package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

func CheckCallbackQuery(call *tg.CallbackQuery) {
	b.TrackAsync(call.From.ID, struct{ *tg.CallbackQuery }{call}, "Callback", func(answer botan.Answer, err []error) {
		log.Ln("Track callback", answer.Status)
		metrika <- true
	})

	usr, err := GetUserDB(call.From.ID)
	if err != nil {
		log.Ln(err.Error())
		return
	}

	log.Ln(usr)

	T, _ := i18n.Tfunc(usr.Language)

	switch {
	case call.Data == "nsfw_true" || call.Data == "nsfw_false":
		ChangeFilter(usr, call, T)
	case call.Data == "settings_menu":
		OpenSettings(usr, call, T)
	case strings.HasPrefix(call.Data, "lang_"):
		switch {
		case call.Data == "lang_menu":
			GetLangList(usr, call, T)
		default:
			ChangeLanguage(usr, call, T)
		}
	}

	<-metrika // Send track to Yandex.metrika
}

func OpenSettings(usr *UserDB, call *tg.CallbackQuery, T i18n.TranslateFunc) {
	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_language"), "lang_menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				T("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(T(fmt.Sprint("status_", usr.NSFW))),
				}),
				fmt.Sprint("nsfw_", !usr.NSFW),
			),
		),
	)

	edit := tg.NewEditMessageReplyMarkup(call.Message.Chat.ID, call.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Ln("Sending message error:", err.Error())
	}
}

func ChangeFilter(usr *UserDB, call *tg.CallbackQuery, T i18n.TranslateFunc) {
	state, _ := strconv.ParseBool(strings.TrimPrefix(call.Data, "nsfw_"))
	go ChangeFilterDB(call.From.ID, state)

	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_language"), "lang_menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				T("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(T(fmt.Sprint("status_", state))),
				}),
				fmt.Sprint("nsfw_", !state),
			),
		),
	)

	edit := tg.NewEditMessageReplyMarkup(call.Message.Chat.ID, call.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Ln("Sending message error:", err.Error())
	}
}

func GetLangList(usr *UserDB, call *tg.CallbackQuery, T i18n.TranslateFunc) {
	var markup tg.InlineKeyboardMarkup
	for _, locale := range locales {
		t, _ := i18n.Tfunc(locale)
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(t("language_name"), fmt.Sprint("lang_", locale)),
		))
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(T("button_cancel"), "settings_menu"),
	))

	edit := tg.NewEditMessageReplyMarkup(call.Message.Chat.ID, call.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Ln("Sending message error:", err.Error())
	}
}

func ChangeLanguage(usr *UserDB, call *tg.CallbackQuery, T i18n.TranslateFunc) {
	newLang := strings.TrimPrefix(call.Data, "lang_")
	go ChangeLangBD(call.From.ID, newLang)

	T, _ = i18n.Tfunc(newLang)

	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_language"), "lang_menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				T("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(T(fmt.Sprint("status_", usr.NSFW))),
				}),
				fmt.Sprint("nsfw_", !usr.NSFW),
			),
		),
	)

	edit := tg.NewEditMessageText(call.Message.Chat.ID, call.Message.MessageID, T("message_settings"))
	edit.ParseMode = parseMarkdown
	edit.ReplyMarkup = &markup
	if _, err := bot.Send(edit); err != nil {
		log.Ln("Sending message error:", err.Error())
	}
}
