package callbacks

import (
	"fmt"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/models"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

const (
	menuLang   = "languages"
	menuRating = "ratings"
	menuRes    = "resources"

	switcherStatus = " 👈🏻"
)

var toggleStatus = map[bool]string{
	true:  "✅ ",
	false: "☑️ ",
}

func CallbackQuery(call *tg.CallbackQuery) {
	usr, err := db.GetUserElseAdd(call.From.ID, call.From.LanguageCode)
	errors.Check(err)

	data := strings.Split(call.Data, ":")
	switch data[0] {
	case "to":
		CallbackTo(usr, call, data[1])
	case "switch":
		callbackSwitch(usr, call, data[1:])
	case "toggle":
		CallbackToggle(usr, call, data[1:])
	case "add":
		callbackAdd(usr, call, data[1:])
	case "remove":
		callbackRemove(usr, call, data[1:])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func CallbackTo(usr *models.User, call *tg.CallbackQuery, option string) {
	log.Ln("CallbackTo", option)
	switch option {
	case "settings":
		CallbackToSettings(usr, call)
	case "languages":
		CallbackToLanguages(usr, call)
	case "resources":
		CallbackToResources(usr, call)
	case "ratings":
		CallbackToRatings(usr, call)
	case "types":
		CallbackToTypes(usr, call)
	case models.BlackList:
		CallbackToList(usr, call, models.BlackList)
	case models.WhiteList:
		CallbackToList(usr, call, models.WhiteList)
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func callbackSwitch(usr *models.User, call *tg.CallbackQuery, options []string) {
	log.Ln("callbackSwitch", options[0])
	switch options[0] {
	case "language":
		CallbackSwitchLanguage(usr, call, options[1])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func CallbackToggle(usr *models.User, call *tg.CallbackQuery, options []string) {
	log.Ln("CallbackToggle", options[0])
	switch options[0] {
	case "resource":
		CallbackToggleResource(usr, call, options[1])
	case "rating":
		CallbackToggleRating(usr, call, options[1])
	case "type":
		CallbackToggleTypes(usr, call, options[1])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}

func callbackAdd(usr *models.User, call *tg.CallbackQuery, options []string) {
	switch options[0] {
	case "tags":
		go func() {
			_, err := bot.Bot.AnswerCallbackQuery(tg.NewAnswerCallbackQuery(call.ID))
			errors.Check(err)
		}()

		usr, err := db.GetUserElseAdd(call.From.ID, call.From.LanguageCode)
		errors.Check(err)

		T, err := i18n.SwitchTo(usr.Language, call.From.LanguageCode)
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

func callbackRemove(usr *models.User, call *tg.CallbackQuery, options []string) {
	switch options[0] {
	case models.BlackList, models.WhiteList:
		err := db.RemoveListTag(usr, options[0], strings.Join(options[1:], ""))
		errors.Check(err)

		CallbackUpdateListKeyboard(usr, call, options[0])
	default:
		CallbackAlert(call, "😱 Oh no!..")
	}
}
