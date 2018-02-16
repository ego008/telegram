package callbacks

import (
	"fmt"

	// log "github.com/kirillDanshin/dlog"
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/models"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	tg "github.com/toby3d/telegram"
)

func CallbackToggleResource(usr *models.User, call *tg.CallbackQuery, res string) {
	err := db.ToggleResource(usr, res)
	errors.Check(err)

	CallbackUpdateResourcesKeyboard(usr, call)
}

func CallbackToResources(usr *models.User, call *tg.CallbackQuery) {
	T, err := i18n.SwitchTo(usr.Language, call.From.LanguageCode)
	errors.Check(err)

	text := T("message_resources")

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = getResourcesMenuKeyboard(usr)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func getResourcesMenuKeyboard(usr *models.User) *tg.InlineKeyboardMarkup {
	T, err := i18n.SwitchTo(usr.Language)
	errors.Check(err)

	var row int
	var replyMarkup tg.InlineKeyboardMarkup
	for i, tag := range resources.Tags {
		if i%2 == 0 {
			replyMarkup.InlineKeyboard = append(
				replyMarkup.InlineKeyboard, tg.NewInlineKeyboardRow(),
			)
			row++
		}

		replyMarkup.InlineKeyboard[row-1] = append(
			replyMarkup.InlineKeyboard[row-1],
			tg.NewInlineKeyboardButton(
				fmt.Sprint(
					toggleStatus[usr.Resources[tag]],
					resources.Resources[tag].UString("title"),
				),
				fmt.Sprint("toggle:resource:", tag),
			),
		)

		i++
	}
	replyMarkup.InlineKeyboard = append(
		replyMarkup.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)

	return &replyMarkup
}

func CallbackUpdateResourcesKeyboard(usr *models.User, call *tg.CallbackQuery) {
	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = getResourcesMenuKeyboard(usr)

	_, err := bot.Bot.EditMessageReplyMarkup(&editMarkup)
	errors.Check(err)
}
