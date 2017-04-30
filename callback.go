package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	i18n "github.com/nicksnyder/go-i18n/i18n"
)

const (
	blacklist  = "blacklist"
	cheatsheet = "cheatsheet"
	menu       = "menu"
	patron     = "patron"
	settings   = "settings"
	whitelist  = "whitelist"
)

var marker = map[bool]string{
	true:  "✅",
	false: "☑️",
}

func callback(callback *tg.CallbackQuery) {
	trackCallback(callback)

	usr, err := getUser(callback.From.ID)
	if err != nil {
		log.Println("Get user:", err.Error())
	}

	T, err := i18n.Tfunc(usr.Language)
	if err != nil {
		log.Println(err.Error())
	}

	data := strings.Split(callback.Data, " ")
	switch data[0] {
	case "info":
		showInfoPopup(usr, callback, T, data)
	case "patreon":
		callPatreon(usr, callback, T, data)
	case "ratings":
		callRatings(usr, callback, T, data)
	case "language":
		callLanguage(usr, callback, T, data)
	case blacklist:
		callTagList(usr, callback, T, data, true)
	case whitelist:
		callTagList(usr, callback, T, data, false)
	case settings:
		callSettings(usr, callback, T, data)
	case "resources":
		callResouces(usr, callback, T, data)
	case "soon":
		call := tg.NewCallback(callback.ID, "Soon™")
		if _, err := bot.AnswerCallbackQuery(call); err != nil {
			log.Println(err.Error())
		}
	}

	<-appMetrika // Send track to Yandex.metrika
}

func callPatreon(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, data []string) {
	switch data[1] {
	case "check":
		pUser, err := p.GetCurrentUser(usr.Patreon.AccessToken)
		if err != nil {
			log.Println(err.Error())
			if _, err := usr.removeRoles(patron); err != nil {
				log.Println(err.Error())
			}
			return
		}

		for _, inc := range pUser.Included {
			if inc.Type == "reward" {
				if inc.Relationships.Campaign.Data.ID == pCampaign {
					if _, err := usr.addRoles(patron); err != nil {
						log.Println(err.Error())
					}

					call := tg.NewCallback(callback.ID, T("message_patreon_has_reward"))
					if _, err := bot.AnswerCallbackQuery(call); err != nil {
						log.Println(err.Error())
					}
					return
				}
			}
		}

		if _, err := usr.removeRoles(patron); err != nil {
			log.Println(err.Error())
		}

		call := tg.NewCallback(callback.ID, T("message_patreon_no_reward"))
		if _, err := bot.AnswerCallbackQuery(call); err != nil {
			log.Println(err.Error())
		}
	case "unlink":
		usr, err = usr.patreonSave("", "", "")
		if err != nil {
			log.Println(err.Error())
		}

		usr, err = usr.removeRoles(patron)
		if err != nil {
			log.Println(err.Error())
		}

		showPatreonKeyboard(usr, callback, T)

		call := tg.NewCallback(callback.ID, T("message_patron_disconnected"))
		if _, err := bot.AnswerCallbackQuery(call); err != nil {
			log.Println(err.Error())
		}
	}
}

func callRatings(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, data []string) {
	switch data[1] {
	case menu:
		showRatingsMessage(usr, callback, T)
		showRatingsKeyboard(usr, callback, T)
	case "change":
		switch data[2] {
		case "s":
			usr, err = usr.changeRatings(!usr.Ratings.Safe, usr.Ratings.Questionale, usr.Ratings.Explicit)
			if err != nil {
				log.Println(err.Error())
			}
			showRatingsKeyboard(usr, callback, T)
		case "q":
			usr, err = usr.changeRatings(usr.Ratings.Safe, !usr.Ratings.Questionale, usr.Ratings.Explicit)
			if err != nil {
				log.Println(err.Error())
			}
			showRatingsKeyboard(usr, callback, T)
		case "e":
			usr, err = usr.changeRatings(usr.Ratings.Safe, usr.Ratings.Questionale, !usr.Ratings.Explicit)
			if err != nil {
				log.Println(err.Error())
			}
			showRatingsKeyboard(usr, callback, T)
		}
	}
}

func callLanguage(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, data []string) {
	switch data[1] {
	case menu:
		showLanguagesMessage(usr, callback, T)
		showLanguagesKeyboard(usr, callback, T)
	case "change":
		usr, err = usr.changeLanguage(data[2])
		if err != nil {
			log.Println(err.Error())
		}

		T, _ = i18n.Tfunc(usr.Language)

		showSettingsMessage(usr, callback, T)
		showSettingsKeyboard(usr, callback, T)

		call := tg.NewCallback(callback.ID, T("message_language_selected"))
		if _, err := bot.AnswerCallbackQuery(call); err != nil {
			log.Println(err.Error())
		}
	}
}

func callTagList(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, data []string, black bool) {
	switch data[1] {
	case menu:
		showTagListMessage(usr, callback, T, black)
		showTagListKeyboard(usr, callback, T, black)
	case "rewrite":
		showTagRewriteMessage(usr, callback, T, black)
	case "remove":
		usr, err = usr.tagRemove(black, data[2])
		if err != nil {
			log.Println(err.Error())
		}
		showTagListKeyboard(usr, callback, T, black)
	}
}

func callSettings(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, data []string) {
	if data[1] == menu {
		showSettingsMessage(usr, callback, T)
		showSettingsKeyboard(usr, callback, T)
	}
}

func callResouces(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, data []string) {
	if data[1] == menu {
		showResourcesMessage(usr, callback, T)
		showResourcesKeyboard(usr, callback, T)
	}
}

func showInfoPopup(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, data []string) {
	rawID := data[1]
	id, _ := strconv.Atoi(rawID[4:])
	posts, err := getPosts(&request{ID: id})
	if err != nil {
		call := tg.NewCallback(callback.ID, "🤷🏻‍♀️")
		if _, err := bot.AnswerCallbackQuery(call); err != nil {
			log.Println(err.Error())
		}
		return
	}
	post := posts[0]

	switch post.Rating {
	case "s":
		post.Rating = strings.Title(T("rating_safe"))
	case "q":
		post.Rating = strings.Title(T("rating_questionable"))
	case "e":
		post.Rating = strings.Title(T("rating_explicit"))
	default:
		post.Rating = strings.Title(T("rating_unknown"))
	}

	text := T("message_post", map[string]interface{}{
		"ID":     post.ID,
		"Posted": time.Unix(post.Change, 0).UTC().Format("2006-01-02 15:04:05"),
		"Owner":  post.Owner,
		"Size":   fmt.Sprint(post.Width, "x", post.Height),
		"Rating": post.Rating,
		"Score":  post.Score,
		"Tags":   post.Tags,
	})
	if len(text) > 197 {
		text = fmt.Sprint(text[:197], "...")
	}
	call := tg.NewCallbackWithAlert(callback.ID, text)
	call.CacheTime = *flagCache
	if _, err := bot.AnswerCallbackQuery(call); err != nil {
		log.Println(err.Error())
	}
}

func showRatingsMessage(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	edit := tg.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, T("message_ratings", map[string]interface{}{
		"Safe":         strings.Title(T("rating_safe")),
		"Questionable": strings.Title(T("rating_questionable")),
		"Explicit":     strings.Title(T("rating_explicit")),
	}))
	edit.ParseMode = tg.ModeMarkdown
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showRatingsKeyboard(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", marker[usr.Ratings.Safe], strings.Title(T("rating_safe"))),
				"ratings change s",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", marker[usr.Ratings.Questionale], strings.Title(T("rating_questionable"))),
				"ratings change q",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", marker[usr.Ratings.Explicit], strings.Title(T("rating_explicit"))),
				"ratings change e",
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_return"), "settings menu"),
		),
	)

	edit := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showLanguagesMessage(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	edit := tg.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, T("message_language"))
	edit.ParseMode = tg.ModeMarkdown
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showLanguagesKeyboard(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	var markup tg.InlineKeyboardMarkup
	for _, locale := range locales {
		t, err := i18n.Tfunc(locale)
		if err != nil {
			log.Println(err.Error())
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", t("language_flag"), t("language_name")),
				fmt.Sprintf("language change %s", locale)),
		))
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(T("button_cancel"), "settings menu"),
	))

	edit := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showSettingsMessage(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	var ratings string
	rt := usr.Ratings
	switch {
	case !rt.Safe && !rt.Questionale && !rt.Explicit, rt.Safe && rt.Questionale && rt.Explicit:
		ratings = T("rating_all")
	case rt.Safe && !rt.Questionale && !rt.Explicit:
		ratings = T("rating_safe")
	case !rt.Safe && rt.Questionale && !rt.Explicit:
		ratings = T("rating_questionable")
	case !rt.Safe && !rt.Questionale && rt.Explicit:
		ratings = T("rating_explicit")
	case !rt.Safe && rt.Questionale && rt.Explicit:
		ratings = fmt.Sprint(T("rating_questionable"), "+", T("rating_explicit"))
	case rt.Safe && !rt.Questionale && rt.Explicit:
		ratings = fmt.Sprint(T("rating_safe"), "+", T("rating_explicit"))
	case rt.Safe && rt.Questionale && !rt.Explicit:
		ratings = fmt.Sprint(T("rating_safe"), "+", T("rating_questionable"))
	}

	text := T("message_settings", map[string]interface{}{
		"Language":  T("language_name"),
		"Ratings":   ratings,
		"Blacklist": strings.Join(usr.Blacklist, " "),
		"Whitelist": strings.Join(usr.Whitelist, " "),
	})
	edit := tg.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)
	edit.ParseMode = tg.ModeMarkdown
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showSettingsKeyboard(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_language", map[string]interface{}{"Flag": T("language_flag")}), "language menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_resources"), "resources menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_ratings"), "ratings menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_blacklist"), "blacklist menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_whitelist"), "whitelist menu"),
		),
	)

	edit := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showPatreonKeyboard(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	pCheckout := &url.URL{
		Scheme: "https",
		Host:   "patreon.com",
		Path:   "/bePatron",
	}
	q := pCheckout.Query()
	q.Add("u", strconv.Itoa(pCampaign))
	pCheckout.RawQuery = q.Encode()

	var markup tg.InlineKeyboardMarkup
	markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonURL(T("button_patreon_checkout"), pCheckout.String()),
	))

	if usr.Patreon.AccessToken == "" {
		pConnect := &url.URL{
			Scheme: "https",
			Host:   "patreon.com",
			Path:   "/oauth2/authorize",
		}
		q := pConnect.Query()
		q.Add("response_type", "code")
		q.Add("client_id", p.ID)
		q.Add("redirect_uri", p.RedirectURI)
		q.Add("scope", "users pledges-to-me")
		q.Add("state", strconv.Itoa(callback.From.ID))
		pConnect.RawQuery = q.Encode()

		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(T("button_patreon_connect"), pConnect.String()),
		))
	} else {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_patreon_check"), "patreon check"),
		))
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				T("button_patreon_disconnect", map[string]interface{}{"Account": usr.Patreon.FullName}),
				"patreon unlink",
			),
		))
	}

	edit := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showTagRewriteMessage(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, black bool) {
	var tagsType string
	if black {
		tagsType = blacklist
	} else {
		tagsType = whitelist
	}

	call := tg.NewCallback(callback.ID, "tag1 tag2 tag* tag_4...")
	if _, err := bot.AnswerCallbackQuery(call); err != nil {
		log.Println(err.Error())
	}

	limit := 5
	for _, role := range usr.Roles {
		if role == patron {
			limit = 15
		}
	}

	reply := tg.NewMessage(
		callback.Message.Chat.ID,
		T(fmt.Sprintf("message_input_%s_tags", tagsType), map[string]interface{}{"Limit": limit}),
	)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = &tg.ForceReply{ForceReply: true}
	if _, err := bot.Send(reply); err != nil {
		log.Println("Sending message error:", err.Error())
	}
}

func showTagListMessage(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, black bool) {
	var text string
	if black {
		text = T("message_blacklist")
	} else {
		text = T("message_whitelist")
	}
	edit := tg.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, text)
	edit.ParseMode = tg.ModeMarkdown
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showTagListKeyboard(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc, black bool) {
	var tags []string
	var list string
	if black {
		list = blacklist
		tags = usr.Blacklist
	} else {
		list = whitelist
		tags = usr.Whitelist
	}

	var markup tg.InlineKeyboardMarkup
	for _, tag := range tags {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(tag, fmt.Sprintf("%s remove %s", list, tag)),
		))
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(T("button_tags_add"), fmt.Sprintf("%s rewrite", list)),
	))
	markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(T("button_return"), "settings menu"),
	))

	edit := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showResourcesMessage(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	edit := tg.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, T("message_resources"))
	edit.ParseMode = tg.ModeMarkdown
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}

func showResourcesKeyboard(usr *User, callback *tg.CallbackQuery, T i18n.TranslateFunc) {
	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s", marker[true], "Gelbooru"), "soon"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_return"), "settings menu"),
		),
	)

	edit := tg.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, markup)
	if _, err := bot.Send(edit); err != nil {
		log.Println(err.Error())
	}
}
