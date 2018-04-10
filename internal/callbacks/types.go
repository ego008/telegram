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

func CallbackToggleTypes(call *tg.CallbackQuery, resultType string) {
	var err error
	switch resultType {
	case models.TypeImage:
		err = database.ToggleTypeImage(call.From)
	case models.TypeAnimation:
		err = database.ToggleTypeAnimation(call.From)
	case models.TypeVideo:
		err = database.ToggleTypeVideo(call.From)
	}
	errors.Check(err)

	if !user.ContentTypes.Animation &&
		!user.ContentTypes.Image &&
		!user.ContentTypes.Video {
		db.ToggleTypeImage(call.From)
		db.ToggleTypeAnimation(call.From)
		db.ToggleTypeVideo(call.From)
	}

	CallbackUpdateTypesKeyboard(call)
}

func CallbackToTypes(call *tg.CallbackQuery) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale, call.From.LanguageCode)
	errors.Check(err)

	editText := tg.NewMessageText(T("message_types"))
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = GetTypesMenuKeyboard(usr)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetTypesMenuKeyboard(user *models.User) *tg.InlineKeyboardMarkup {
	T, err := i18n.SwitchTo(user.Locale)
	errors.Check(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[user.ContentTypes.Image],
					strings.Title(T("type_image")),
				),
				"toggle:type:image",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[user.ContentTypes.Animation],
					strings.Title(T("type_animation")),
				),
				"toggle:type:animation",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[user.ContentTypes.Video],
					strings.Title(T("type_video")),
				),
				"toggle:type:video",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)
}

func CallbackUpdateTypesKeyboard(call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = GetTypesMenuKeyboard(usr)

	_, err := bot.Bot.EditMessageReplyMarkup(&editMarkup)
	errors.Check(err)
}
