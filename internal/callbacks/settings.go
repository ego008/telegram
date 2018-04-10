package callbacks

import (
	"strings"

	// log "github.com/kirillDanshin/dlog"
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func CallbackToSettings(call *tg.CallbackQuery) {
	T, err := i18n.SwitchTo(user.Locale, call.From.LanguageCode)
	errors.Check(err)

	var activeRes []string
	for k, v := range user.Resources {
		if v && resources.Resources[k] != nil {
			activeRes = append(activeRes, resources.Resources[k].UString("title"))
		}
	}

	ratings, err := user.GetRatingsStatus()
	errors.Check(err)

	text := T("message_settings", map[string]interface{}{
		"Language":  i18n.Names[user.Locale],
		"Resources": strings.Join(activeRes, "`, `"),
		"Ratings":   ratings,
		"Blacklist": strings.Join(user.Blacklist, "`, `"),
		"Whitelist": strings.Join(user.Whitelist, "`, `"),
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = GetSettingsMenuKeyboard(usr)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetSettingsMenuKeyboard(user *models.User) *tg.InlineKeyboardMarkup {
	T, err := i18n.SwitchTo(user.Locale)
	errors.Check(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_language", map[string]interface{}{
					"Flag": T("language_flag"),
				}),
				"to:languages",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_resources"), "to:resources",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_ratings"), "to:ratings",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_types"), "to:types",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_blacklist"), "to:blacklist",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_whitelist"), "to:whitelist",
			),
		),
	)
}
