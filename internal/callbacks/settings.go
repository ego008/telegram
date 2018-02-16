package callbacks

import (
	"strings"

	// log "github.com/kirillDanshin/dlog"
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/models"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	tg "github.com/toby3d/telegram"
)

func CallbackToSettings(usr *models.User, call *tg.CallbackQuery) {
	T, err := i18n.SwitchTo(usr.Language, call.From.LanguageCode)
	errors.Check(err)

	var activeRes []string
	for k, v := range usr.Resources {
		if v && resources.Resources[k] != nil {
			activeRes = append(activeRes, resources.Resources[k].UString("title"))
		}
	}

	ratings, err := usr.GetRatingsStatus()
	errors.Check(err)

	text := T("message_settings", map[string]interface{}{
		"Language":  i18n.Names[usr.Language],
		"Resources": strings.Join(activeRes, "`, `"),
		"Ratings":   ratings,
		"Blacklist": strings.Join(usr.Blacklist, "`, `"),
		"Whitelist": strings.Join(usr.Whitelist, "`, `"),
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = GetSettingsMenuKeyboard(usr)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetSettingsMenuKeyboard(usr *models.User) *tg.InlineKeyboardMarkup {
	T, err := i18n.SwitchTo(usr.Language)
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
				"Type filters", "to:types",
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
