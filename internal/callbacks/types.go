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

func CallbackToggleTypes(usr *models.User, call *tg.CallbackQuery, resultType string) {
	var err error
	switch resultType {
	case "image":
		err = db.ToggleTypeImage(usr)
	case "animation":
		err = db.ToggleTypeAnimation(usr)
	case "video":
		err = db.ToggleTypeVideo(usr)
	}
	errors.Check(err)

	if !usr.ContentTypes.Animation &&
		!usr.ContentTypes.Image &&
		!usr.ContentTypes.Video {
		db.ToggleTypeImage(usr)
		db.ToggleTypeAnimation(usr)
		db.ToggleTypeVideo(usr)
	}

	CallbackUpdateTypesKeyboard(usr, call)
}

func CallbackToTypes(usr *models.User, call *tg.CallbackQuery) {
	T, err := i18n.SwitchTo(usr.Language, call.From.LanguageCode)
	errors.Check(err)

	editText := tg.NewMessageText(T("message_types"))
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = GetTypesMenuKeyboard(usr)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetTypesMenuKeyboard(usr *models.User) *tg.InlineKeyboardMarkup {
	T, err := i18n.SwitchTo(usr.Language)
	errors.Check(err)

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.ContentTypes.Image],
					strings.Title(T("type_image")),
				),
				"toggle:type:image",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.ContentTypes.Animation],
					strings.Title(T("type_animation")),
				),
				"toggle:type:animation",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.ContentTypes.Video],
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

func CallbackUpdateTypesKeyboard(usr *models.User, call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = GetTypesMenuKeyboard(usr)

	_, err := bot.Bot.EditMessageReplyMarkup(&editMarkup)
	errors.Check(err)
}
