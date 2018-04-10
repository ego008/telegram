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

func CallbackToggleRating(call *tg.CallbackQuery, rating string) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	switch rating {
	case models.RatingSafe:
		err = database.DB.ToggleRatingSafe(call.From)
	case models.RatingQuestionable:
		err = database.DB.ToggleRatingQuestionable(call.From)
	case models.RatingExplicit:
		err = database.DB.ToggleRatingExplicit(call.From)
	}
	errors.Check(err)

	if !user.Ratings.Safe &&
		!user.Ratings.Questionable &&
		!user.Ratings.Exlplicit {
		database.DB.ToggleRatingSafe(call.From)
		database.DB.ToggleRatingQuestionable(call.From)
		database.DB.ToggleRatingExplicit(call.From)
	}

	CallbackUpdateRatingsKeyboard(call)
}

func CallbackToRatings(call *tg.CallbackQuery) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale, call.From.LanguageCode)
	errors.Check(err)

	text := T("message_ratings", map[string]interface{}{
		"Safe":         strings.Title(T("rating_safe")),
		"Questionable": strings.Title(T("rating_questionable")),
		"Explicit":     strings.Title(T("rating_explicit")),
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = GetRatingsMenuKeyboard(call)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetRatingsMenuKeyboard(call *tg.CallbackQuery) *tg.InlineKeyboardMarkup {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale)
	errors.Check(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[user.Ratings.Safe],
					strings.Title(T("rating_safe")),
				),
				"toggle:rating:safe",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[user.Ratings.Questionable],
					strings.Title(T("rating_questionable")),
				),
				"toggle:rating:questionable",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[user.Ratings.Exlplicit],
					strings.Title(T("rating_explicit")),
				),
				"toggle:rating:explicit",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)
}

func CallbackUpdateRatingsKeyboard(call *tg.CallbackQuery) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = GetRatingsMenuKeyboard(call)

	_, err = bot.Bot.EditMessageReplyMarkup(&editMarkup)
	errors.Check(err)
}
