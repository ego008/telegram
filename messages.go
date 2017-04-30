package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	durafmt "github.com/hako/durafmt"
	i18n "github.com/nicksnyder/go-i18n/i18n"
	http "github.com/valyala/fasthttp"
)

var startUptime = time.Now()

func message(msg *tg.Message) {
	if msg.From.ID == bot.Self.ID {
		return
	}

	if !msg.Chat.IsPrivate() {
		if _, err := bot.LeaveChat(tg.ChatConfig{ChatID: msg.Chat.ID}); err != nil {
			log.Println(err.Error())
			return
		}
	}

	usr, err := getUser(msg.From.ID)
	if err != nil {
		trackMessage(msg, "Message")
		log.Println("Get user:", err.Error())
		<-appMetrika
		return
	}

	T, err := i18n.Tfunc(usr.Language)
	if err != nil {
		trackMessage(msg, "Message")
		log.Println(err.Error())
		<-appMetrika
		return
	}

	isCommand := msg.IsCommand()
	isPrivate := msg.Chat.IsPrivate()
	switch {
	case isCommand:
		switch strings.ToLower(msg.Command()) {
		case "start": // Requirement Telegram platform
			cmdStart(usr, msg, T)
		case "help": // Requirement Telegram platform
			cmdHelp(usr, msg, T)
		case "settings":
			cmdSettings(usr, msg, T)
		case "cheatsheet":
			cmdCheatsheet(usr, msg, T)
		case "info":
			cmdInfo(usr, msg, T)
		case "patreon":
			cmdPatreon(usr, msg, T)
		default:
			cmdEasterEgg(msg)
		}
	case !isCommand && msg.ReplyToMessage != nil:
		if msg.ReplyToMessage.Text == "" || msg.Text == "" {
			trackMessage(msg, "Message")
			<-appMetrika
			return
		}
		analyzeReply(usr, msg, T)
	case !isCommand && isPrivate && msg.Text == "":
		for _, role := range usr.Roles {
			if role == "admin" {
				getTelegramFileID(msg) // Admin feature without tracking
			}
		}
	default:
		msgEasterEgg(msg) // Secret actions and commands ;)
	}
}

func analyzeReply(usr *User, msg *tg.Message, T i18n.TranslateFunc) {
	trackMessage(msg, "/settings") // Track action

	limit := 5
	for _, role := range usr.Roles {
		if role == patron {
			limit = 15
		}
	}

	switch msg.ReplyToMessage.Text {
	case T("message_input_blacklist_tags", map[string]interface{}{"Limit": limit}):
		tags := strings.Split(strings.ToLower(msg.Text), " ")

		limit := 5
		for _, role := range usr.Roles {
			if role == patron {
				limit = 15
			}
		}

		if len(tags) >= limit {
			tags = tags[:limit]
		}

		if err := usr.tagsRewrite(true, tags); err != nil {
			log.Println(err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}

		reply := tg.NewMessage(
			msg.Chat.ID,
			T(
				"message_input_tags_success",
				map[string]interface{}{
					"Tags": strings.Join(tags, " "),
				},
			),
		)
		reply.ParseMode = tg.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			log.Println("Sending message error:", err.Error())
		}
	case T("message_input_whitelist_tags", map[string]interface{}{"Limit": limit}):
		tags := strings.Split(strings.ToLower(msg.Text), " ")

		limit := 5
		for _, role := range usr.Roles {
			if role == patron {
				limit = 15
			}
		}

		if len(tags) >= limit {
			tags = tags[:limit]
		}

		if err := usr.tagsRewrite(false, tags); err != nil {
			log.Println(err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}

		reply := tg.NewMessage(
			msg.Chat.ID,
			T(
				"message_input_tags_success",
				map[string]interface{}{
					"Tags": strings.Join(tags, " "),
				},
			),
		)
		reply.ParseMode = tg.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			log.Println(err.Error())
		}
	}
	<-appMetrika // Send track to Yandex.metrika
}

func cmdStart(usr *User, msg *tg.Message, T i18n.TranslateFunc) {
	trackMessage(msg, "/start") // Track action

	args := strings.Split(msg.CommandArguments(), " ")
	switch {
	case args[0] == "settings":
		cmdSettings(usr, msg, T)
		return
	case args[0] == "cheatsheet":
		cmdCheatsheet(usr, msg, T)
		return
	case strings.HasPrefix(args[0], "code_"):
		args = strings.Split(args[0], "_")
		id, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println(err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}

		if id != msg.From.ID {
			log.Println("not contain or not compare id")
			<-appMetrika // Send track to Yandex.metrika
			return
		}

		resp, err := p.ValidateReceipt(args[2])
		if err != nil {
			log.Println(err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}

		pUser, err := p.GetCurrentUser(resp.AccessToken)
		if err != nil {
			log.Println(err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}

		usr, err = usr.patreonSave(pUser.Data.Attributes.FullName, resp.AccessToken, resp.RefreshToken)
		if err != nil {
			log.Println(err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}

		reply := tg.NewMessage(
			msg.Chat.ID,
			T(
				"message_patron_connected",
				map[string]interface{}{
					"FullName": usr.Patreon.FullName,
				},
			),
		)
		reply.ParseMode = tg.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			log.Println("Sending message error:", err.Error())
		}

		<-appMetrika // Send track to Yandex.metrika
		return
	}

	bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)) // Force feedback

	demo := cfg.UString("telegram.content.demo")
	if demo != "" {
		document := tg.NewDocumentShare(msg.Chat.ID, demo)
		if _, err := bot.Send(document); err != nil {
			log.Println("Sending message error:", err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}
	}

	exampleQuery := "hatsune_miku rating:safe"
	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.InlineKeyboardButton{
				Text: T("button_try"),
				SwitchInlineQueryCurrentChat: &exampleQuery,
			},
		),
	)

	text := T("message_start", map[string]interface{}{"FirstName": msg.From.FirstName, "BotName": bot.Self.UserName})
	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = &markup
	if _, err := bot.Send(reply); err != nil {
		log.Println("Sending message error:", err.Error())
	}

	<-appMetrika // Send track to Yandex.metrika
}

func cmdHelp(usr *User, msg *tg.Message, T i18n.TranslateFunc) {
	trackMessage(msg, "/help")                             // Track action
	bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)) // Force feedback

	demo := cfg.UString("telegram.content.demo")
	if demo != "" {
		document := tg.NewDocumentShare(msg.Chat.ID, demo)
		if _, err := bot.Send(document); err != nil {
			log.Println("Sending message error:", err.Error())
			<-appMetrika // Send track to Yandex.metrika
			return
		}
	}

	reply := tg.NewMessage(msg.Chat.ID, T("message_help"))
	reply.ParseMode = tg.ModeMarkdown
	reply.DisableWebPagePreview = true
	if _, err := bot.Send(reply); err != nil {
		log.Println("Sending message error:", err.Error())
	}

	<-appMetrika // Send track to Yandex.metrika
}

func cmdSettings(usr *User, msg *tg.Message, T i18n.TranslateFunc) {
	trackMessage(msg, "/settings")                         // Track action
	bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)) // Force feedback

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

	text := T("message_settings", map[string]interface{}{
		"Language":  T("language_name"),
		"Ratings":   ratings,
		"Blacklist": strings.Join(usr.Blacklist, " "),
		"Whitelist": strings.Join(usr.Whitelist, " "),
	})
	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyMarkup = &markup
	if _, err := bot.Send(reply); err != nil {
		log.Println("Sending message error:", err.Error())
	}

	<-appMetrika // Send track to Yandex.metrika
}

func cmdCheatsheet(usr *User, msg *tg.Message, T i18n.TranslateFunc) {
	trackMessage(msg, "/cheatsheet")                       // Track action
	bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)) // Force feedback

	reply := tg.NewMessage(msg.Chat.ID, T("message_cheatsheet"))
	reply.ParseMode = tg.ModeMarkdown
	reply.DisableWebPagePreview = true
	if _, err := bot.Send(reply); err != nil {
		log.Println("Sending message error:", err.Error())
	}

	<-appMetrika // Send track to Yandex.metrika
}

func cmdPatreon(usr *User, msg *tg.Message, T i18n.TranslateFunc) {
	trackMessage(msg, "/patreon")                          // Track action
	bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)) // Force feedback

	users, err := getUsers()
	if err != nil {
		log.Println(err.Error())
		<-appMetrika // Send track to Yandex.metrika
		return
	}

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
		q.Add("state", strconv.Itoa(msg.From.ID))
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

	text := T("message_patreon_empty")
	var names []string
	for _, user := range users {
		isPatron := false
		for _, role := range user.Roles {
			if role == patron {
				isPatron = true
			}
		}
		if isPatron && user.Patreon.FullName != "" {
			names = append(names, user.Patreon.FullName)
		}
	}

	if len(names) > 0 {
		text = T("message_patreon", map[string]interface{}{
			"Patrons": strings.Join(names, "\n"),
		})
	}

	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = tg.ModeMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyMarkup = &markup
	if _, err := bot.Send(reply); err != nil {
		log.Println("Sending message error:", err.Error())
	}

	<-appMetrika // Send track to Yandex.metrika
}

func cmdInfo(usr *User, msg *tg.Message, T i18n.TranslateFunc) {
	trackMessage(msg, "/info")                             // Track action
	bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)) // Force feedback

	uptimePeriod := time.Since(startUptime).String()
	uptime, err := durafmt.ParseString(uptimePeriod)
	if err != nil {
		log.Println(err.Error())
		<-appMetrika // Send track to Yandex.metrika
		return
	}

	var markup tg.InlineKeyboardMarkup
	var social []tg.InlineKeyboardButton

	chInvite := cfg.UString("telegram.channel.invite")
	if chInvite != "" {
		social = append(social, tg.NewInlineKeyboardButtonURL(
			T("button_channel"), chInvite,
		))
	}

	gInvite := cfg.UString("telegram.group.invite")
	if gInvite != "" {
		social = append(social, tg.NewInlineKeyboardButtonURL(
			T("button_group"), gInvite,
		))
	}

	dInvite := cfg.UString("discord.invite")
	if dInvite != "" {
		social = append(social, tg.NewInlineKeyboardButtonURL(
			T("button_discord"), dInvite,
		))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, social)
	markup.InlineKeyboard = append(markup.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonURL(T("button_rate"), fmt.Sprint("https://t.me/storebot?start=", bot.Self.UserName)),
	))

	cover := cfg.UString("telegram.content.cover")
	if cover != "" {
		photo := tg.NewPhotoShare(msg.Chat.ID, cover)
		photo.Caption = T("message_info_caption", map[string]interface{}{
			"Version": fmt.Sprintf("%s (%d)", version, build),
			"UpTime":  uptime.String(),
		})
		photo.ReplyMarkup = &markup
		if _, err := bot.Send(photo); err != nil {
			log.Println("Sending message error:", err.Error())
		}
	} else {
		reply := tg.NewMessage(msg.Chat.ID, T("message_info", map[string]interface{}{
			"Version": fmt.Sprintf("%s (%d)", version, build),
			"UpTime":  uptime.String(),
		}))
		reply.ReplyMarkup = &markup
		reply.ParseMode = tg.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			log.Println("Sending message error:", err.Error())
		}
	}

	<-appMetrika // Send track to Yandex.metrika
}

func getTelegramFileID(msg *tg.Message) {
	var uploadFileInfo string
	switch {
	case msg.Audio != nil: // Upload file as Voice
		if strings.Contains(msg.Audio.MimeType, "ogg") {
			voice, err := getVoiceFromAudio(msg)
			if err != nil {
				log.Println(err.Error())
				return
			}
			uploadFileInfo = fmt.Sprintf("ID: %s", voice)
		} else {
			uploadFileInfo = fmt.Sprintf("ID: %s", msg.Audio.FileID)
		}
	case msg.Document != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", msg.Document.FileID)
	case msg.Photo != nil: // Get large file ID
		photos := *msg.Photo
		uploadFileInfo = fmt.Sprintf("ID: %s", photos[len(photos)-1].FileID)
	case msg.Sticker != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", msg.Sticker.FileID)
	case msg.Video != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", msg.Video.FileID)
	case msg.Voice != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", msg.Voice.FileID)
	}
	reply := tg.NewMessage(msg.Chat.ID, uploadFileInfo)
	reply.ReplyToMessageID = msg.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Println(err.Error())
	}
}

func getVoiceFromAudio(msg *tg.Message) (string, error) {
	bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatRecordAudio))

	link, err := bot.GetFileDirectURL(msg.Audio.FileID)
	if err != nil {
		return "", err
	}

	_, body, err := http.Get(nil, link)
	if err != nil {
		return "", err
	}
	bytes := tg.FileBytes{
		Name:  msg.Audio.FileID,
		Bytes: body,
	}

	voice := tg.NewVoiceUpload(msg.Chat.ID, bytes)
	voice.Duration = msg.Audio.Duration
	voice.ReplyToMessageID = msg.MessageID
	resp, err := bot.Send(voice)
	if err != nil {
		return "", err
	}

	return resp.Voice.FileID, nil
}
