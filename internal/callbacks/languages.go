package callbacks

import (
	"fmt"
	"strings"

	// log "github.com/kirillDanshin/dlog"
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func CallbackSwitchLanguage(call *tg.CallbackQuery, lang string) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	if lang == user.Locale {
		// Because we must answer on every callback request
		_, err := bot.Bot.AnswerCallbackQuery(
			tg.NewAnswerCallbackQuery(call.ID),
		)
		errors.Check(err)
		return
	}

	err = database.DB.SetLocale(call.From, lang)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale, call.From.LanguageCode)
	errors.Check(err)

	go CallbackAlert(call, T("message_language_selected"))

	CallbackToLanguages(call)
}

func CallbackToLanguages(call *tg.CallbackQuery) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale, call.From.LanguageCode)
	errors.Check(err)

	text := T("message_language", map[string]interface{}{
		"LanguageCodes": strings.Join(i18n.Tags, "|"),
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = GetLanguagesMenuKeyboard(user)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetLanguagesMenuKeyboard(user *models.User) *tg.InlineKeyboardMarkup {
	T, err := i18n.SwitchTo(user.Locale)
	errors.Check(err)

	var replyMarkup tg.InlineKeyboardMarkup
	for _, tag := range i18n.Tags {
		var this string
		if user.Locale == tag {
			this = switcherStatus
		}

		replyMarkup.InlineKeyboard = append(
			replyMarkup.InlineKeyboard,
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButton(
					fmt.Sprint(i18n.Names[tag], this),
					fmt.Sprint("switch:language:", tag),
				),
			),
		)
	}
	replyMarkup.InlineKeyboard = append(
		replyMarkup.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)

	return &replyMarkup
}

func callbackUpdateLanguagesKeyboard(call *tg.CallbackQuery) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = GetLanguagesMenuKeyboard(user)

	_, err = bot.Bot.EditMessageReplyMarkup(&editMarkup)
	errors.Check(err)
}
