package callbacks

import (
	"fmt"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

const switcherStatus = " 👈🏻"

var toggleStatus = map[bool]string{
	true:  "✅ ",
	false: "☑️ ",
}

func CallbackQuery(call *tg.CallbackQuery) {
	data := strings.Split(call.Data, ":")
	switch data[0] {
	case "to":
		CallbackTo(call, data[1])
	case "switch":
		callbackSwitch(call, data[1:])
	case "toggle":
		CallbackToggle(call, data[1:])
	case "add":
		callbackAdd(call, data[1:])
	case "remove":
		callbackRemove(call, data[1:])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func CallbackTo(call *tg.CallbackQuery, option string) {
	log.Ln("CallbackTo", option)
	switch option {
	case "settings":
		CallbackToSettings(call)
	case "languages":
		CallbackToLanguages(call)
	case "resources":
		CallbackToResources(call)
	case "ratings":
		CallbackToRatings(call)
	case "types":
		CallbackToTypes(call)
	case models.BlackList:
		CallbackToList(call, models.BlackList)
	case models.WhiteList:
		CallbackToList(call, models.WhiteList)
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func callbackSwitch(call *tg.CallbackQuery, options []string) {
	log.Ln("callbackSwitch", options[0])
	switch options[0] {
	case "language":
		CallbackSwitchLanguage(call, options[1])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func CallbackToggle(call *tg.CallbackQuery, options []string) {
	log.Ln("CallbackToggle", options[0])
	switch options[0] {
	case "resource":
		CallbackToggleResource(call, options[1])
	case "rating":
		CallbackToggleRating(call, options[1])
	case "type":
		CallbackToggleTypes(call, options[1])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func callbackAdd(call *tg.CallbackQuery, options []string) {
	user, err := database.DB.GetUser(call.From)
	errors.Check(err)

	switch options[0] {
	case "tags":
		go func() {
			_, err := bot.Bot.AnswerCallbackQuery(tg.NewAnswerCallbackQuery(call.ID))
			errors.Check(err)
		}()

		usr, err := database.DB.GetUser(call.From)
		errors.Check(err)

		T, err := i18n.SwitchTo(user.Locale, call.From.LanguageCode)
		errors.Check(err)

		text := T(fmt.Sprint("message_input_", options[1], "_tags"), map[string]interface{}{
			"Limit": 5,
		})

		reply := tg.NewMessage(int64(call.From.ID), text)
		reply.ParseMode = tg.ModeMarkdown
		reply.ReplyMarkup = tg.NewForceReply(true)

		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func callbackRemove(call *tg.CallbackQuery, options []string) {
	switch options[0] {
	case models.BlackList, models.WhiteList:
		var err error
		switch options[0] {
		case models.BlackList:
			err = database.DB.RemoveBlackTag(call.From, strings.Join(options[1:], ""))
		case models.WhiteList:
			err = database.DB.RemoveWhiteTag(call.From, strings.Join(options[1:], ""))
		}
		errors.Check(err)

		CallbackUpdateListKeyboard(call, options[0])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}
