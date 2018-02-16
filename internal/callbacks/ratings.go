package callbacks

import (
	"fmt"
	"strings"

	// log "github.com/kirillDanshin/dlog"
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/models"
	tg "github.com/toby3d/telegram"
)

func CallbackToggleRating(usr *models.User, call *tg.CallbackQuery, rating string) {
	var err error
	switch rating {
	case "safe":
		err = db.ToggleRatingSafe(usr)
	case "questionable":
		err = db.ToggleRatingQuestionable(usr)
	case "explicit":
		err = db.ToggleRatingExplicit(usr)
	}
	errors.Check(err)

	if !usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit {
		db.ToggleRatingSafe(usr)
		db.ToggleRatingQuestionable(usr)
		db.ToggleRatingExplicit(usr)
	}

	CallbackUpdateRatingsKeyboard(usr, call)
}

func CallbackToRatings(usr *models.User, call *tg.CallbackQuery) {
	T, err := i18n.SwitchTo(usr.Language, call.From.LanguageCode)
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
	editText.ReplyMarkup = GetRatingsMenuKeyboard(usr)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetRatingsMenuKeyboard(usr *models.User) *tg.InlineKeyboardMarkup {
	T, err := i18n.SwitchTo(usr.Language)
	errors.Check(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Ratings.Safe],
					strings.Title(T("rating_safe")),
				),
				"toggle:rating:safe",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Ratings.Questionable],
					strings.Title(T("rating_questionable")),
				),
				"toggle:rating:questionable",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Ratings.Exlplicit],
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

func CallbackUpdateRatingsKeyboard(usr *models.User, call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = GetRatingsMenuKeyboard(usr)

	_, err := bot.Bot.EditMessageReplyMarkup(&editMarkup)
	errors.Check(err)
}
