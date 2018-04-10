package callbacks

import (
	"fmt"

	// log "github.com/kirillDanshin/dlog"
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func CallbackToList(call *tg.CallbackQuery, listType string) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale, call.From.LanguageCode)
	errors.Check(err)

	text := T(fmt.Sprint("message_", listType), map[string]interface{}{
		"CommandCheatsheet": models.CommandCheatsheet,
	})

	editText := tg.NewMessageText(text)
	editText.ChatID = call.Message.Chat.ID
	editText.MessageID = call.Message.ID
	editText.ParseMode = tg.ModeMarkdown
	editText.ReplyMarkup = GetListMenuKeyboard(call, listType)

	_, err = bot.Bot.EditMessageText(editText)
	errors.Check(err)
}

func GetListMenuKeyboard(call *tg.CallbackQuery, listType string) *tg.InlineKeyboardMarkup {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	T, err := i18n.SwitchTo(user.Locale)
	errors.Check(err)

	var tags []string
	switch listType {
	case models.BlackList:
		tags = user.BlackTags()
	case models.WhiteList:
		tags = user.WhiteTags()
	}

	replyMarkup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(
				T("button_tags_add"), fmt.Sprint("add:tags:", listType),
			),
		),
	)
	row := 1
	for i, tag := range tags {
		if i%2 == 0 {
			replyMarkup.InlineKeyboard = append(
				replyMarkup.InlineKeyboard, tg.NewInlineKeyboardRow(),
			)
			row++
		}

		replyMarkup.InlineKeyboard[row-1] = append(
			replyMarkup.InlineKeyboard[row-1],
			tg.NewInlineKeyboardButton(
				tag, fmt.Sprint("remove:", listType, ":", tag),
			),
		)
	}
	replyMarkup.InlineKeyboard = append(
		replyMarkup.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButton(T("button_return"), "to:settings"),
		),
	)

	return replyMarkup
}

func CallbackUpdateListKeyboard(call *tg.CallbackQuery, listType string) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	var editMarkup tg.EditMessageReplyMarkupParameters
	editMarkup.ChatID = call.Message.Chat.ID
	editMarkup.MessageID = call.Message.ID
	editMarkup.ReplyMarkup = GetListMenuKeyboard(call, listType)

	_, err = bot.Bot.EditMessageReplyMarkup(&editMarkup)
	errors.Check(err)
}
