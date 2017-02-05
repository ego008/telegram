package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/botanio/sdk/go"
	"github.com/hako/durafmt"
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
	f "github.com/valyala/fasthttp"
	tg "gopkg.in/telegram-bot-api.v4"
)

const (
	parseMarkdown    = "markdown"
	parseHTML        = "html"
	versionCover     = "AgADAgAD0qcxG5dBuEiCC-cbcKOJf_U2Sw0ABLXb7jf__uQVOr4EAAEC"
	versionCodeName  = "3.0 \"Cinder Fall\""
	demonstrationGIF = "BQADAgADNwYAAmOYSwLFYMl_HVAaDwI"
)

var startUptime = time.Now()

func GetMessage(msg *tg.Message) {
	usr, err := GetUserDB(msg.From.ID)
	if err != nil {
		log.Ln(err.Error())
		return
	}

	T, _ := i18n.Tfunc(usr.Language)

	if usr.Role == "anon" {
		markup := tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData(T("button_verify"), "i_agree"),
			),
		)
		reply := tg.NewMessage(int64(msg.From.ID), T("message_verify"))
		reply.ParseMode = parseMarkdown
		reply.DisableWebPagePreview = true
		reply.ReplyMarkup = &markup
		if _, err := bot.Send(reply); err != nil {
			log.Ln("Sending message error:", err.Error())
		}
		return
	}

	isCommand := msg.IsCommand()
	isPrivate := msg.Chat.IsPrivate()
	switch {
	case isCommand:
		Commands(usr, msg, T)
	case !isCommand && isPrivate && usr.Role == "admin" && msg.Text == "":
		getTelegramFileID(msg) // Admin feature without tracking
	default:
		EasterEggsMessages(msg) // Secret actions and commands ;)
	}
}

func Commands(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	lowerCommand := strings.ToLower(msg.Command())
	switch lowerCommand {
	case "start": // Requirement Telegram platform
		StartCommand(usr, msg, T)
	case "help": // Requirement Telegram platform
		HelpCommand(usr, msg, T)
	case "settings": // Requirement Telegram platform
		go SettingsCommand(usr, msg, T)
	case "cheatsheet":
		CheatSheetCommand(usr, msg, T)
	case "random":
		RandomCommand(usr, msg, T)
	case "info":
		InfoCommand(usr, msg, T)
	case "donate":
		DonateCommand(usr, msg, T)
	default:
		EggCommand(msg)
	}
}

func StartCommand(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(msg.From.ID, struct{ *tg.Message }{msg}, "/start", func(answer botan.Answer, err []error) {
		log.Ln("Track /start %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	switch msg.CommandArguments() {
	case "settings":
		SettingsCommand(usr, msg, T)
		return
	}

	inlineKeyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitch(T("button_try"), "hatsune_miku rating:safe"), // Showing tutorial button for demonstration work
		),
	)

	text := T("message_start", map[string]interface{}{"FirstName": msg.From.FirstName})
	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Ln("Sending message error:", err.Error())
	}

	<-metrika // Send track to Yandex.metrika
}

func HelpCommand(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(msg.From.ID, struct{ *tg.Message }{msg}, "/help", func(answer botan.Answer, err []error) {
		log.Ln("Track /help", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	document := tg.NewDocumentShare(int64(msg.From.ID), demonstrationGIF)
	if _, err := bot.Send(document); err != nil {
		log.Ln("Sending message error:", err.Error())
	}

	text := T("message_help")
	reply := tg.NewMessage(int64(msg.From.ID), text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	if _, err := bot.Send(reply); err != nil {
		log.Ln("Sending message error:", err.Error())
	}

	<-metrika // Send track to Yandex.metrika
}

func SettingsCommand(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(msg.From.ID, struct{ *tg.Message }{msg}, "/settings", func(answer botan.Answer, err []error) {
		log.Ln("Track /settings", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatTyping)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_language"), "lang_menu"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				T("button_nsfw", map[string]interface{}{
					"Status": strings.ToUpper(T(fmt.Sprint("status_", usr.NSFW))),
				}),
				fmt.Sprint("nsfw_", !usr.NSFW),
			),
		),
	)

	text := T("message_settings")
	reply := tg.NewMessage(int64(msg.From.ID), text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyMarkup = &markup
	if _, err := bot.Send(reply); err != nil {
		log.Ln("Sending message error:", err.Error())
	}

	<-metrika // Send track to Yandex.metrika
}

func CheatSheetCommand(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(msg.From.ID, struct{ *tg.Message }{msg}, "/cheatsheet", func(answer botan.Answer, err []error) {
		log.Ln("Track /cheatsheet", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(int64(msg.From.ID), tg.ChatTyping)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	text := T("message_cheatsheet")
	reply := tg.NewMessage(int64(msg.From.ID), text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	if _, err := bot.Send(reply); err != nil {
		log.Ln("Sending message error:", err.Error())
	}

	<-metrika // Send track to Yandex.metrika
}

func DonateCommand(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(msg.From.ID, struct{ *tg.Message }{msg}, "/donate", func(answer botan.Answer, err []error) {
		log.Ln("Track /donate %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(int64(msg.From.ID), tg.ChatTyping)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	var donateURL string
	if msg.Chat.IsPrivate() {
		donateURL = getBotanURL(msg.From.ID, cfg["link_donate"].(string))
	} else {
		donateURL = cfg["link_donate"].(string)
	}

	inlineKeyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(T("button_donate"), donateURL),
		),
	)

	text := T("message_donate")
	reply := tg.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Ln("Sending message error:", err.Error())
	}

	<-metrika // Send track to Yandex.metrika
}

func RandomCommand(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(msg.From.ID, struct{ *tg.Message }{msg}, "/random", func(answer botan.Answer, err []error) {
		log.Ln("Track /random", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatUploadDocument)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	randomSource := rand.NewSource(time.Now().UnixNano()) // Maximum randomizing dice
	totalPosts := getPosts(Request{ID: 0})                // Get last upload post
	random := rand.New(randomSource)                      // Create magical dice
	var randomFile []Post

	for {
		randomPost := random.Intn(totalPosts[0].ID)    // Generate a random ID number from first to last ID post
		randomFile = getPosts(Request{ID: randomPost}) // Call to selected ID
		if len(randomFile) > 0 {
			break // If post is NOT blocked or erroneous
		}
		log.Ln("Nothing. Reroll dice!")
	}

	log.Ln("Get random file by URL: %s", randomFile[0].FileURL)
	_, body, err := f.Get(nil, randomFile[0].FileURL)
	if err != nil {
		log.Ln("Get random image by URL error:", err.Error())
	}
	bytes := tg.FileBytes{
		Name:  randomFile[0].Image,
		Bytes: body,
	}
	uploadFilesProcess(msg, bytes, randomFile[0], T)

	<-metrika // Send track to Yandex.metrika
}

func InfoCommand(usr *UserDB, msg *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(msg.From.ID, struct{ *tg.Message }{msg}, "/info", func(answer botan.Answer, err []error) {
		log.Ln("Track /info", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(int64(msg.From.ID), tg.ChatTyping)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	uptimePeriod := time.Since(startUptime).String()
	uptime, err := durafmt.ParseString(uptimePeriod)
	if err != nil {
		fmt.Printf("DuraFmt error:", err.Error())
	}

	inlineKeyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(T("button_channel"), cfg["link_channel"].(string)),
			tg.NewInlineKeyboardButtonURL(T("button_group"), cfg["link_group"].(string)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(T("button_rate"), cfg["link_rate"].(string)+bot.Self.UserName),
		),
	)
	photo := tg.NewPhotoShare(int64(msg.From.ID), versionCover)
	photo.Caption = T("message_info", map[string]interface{}{
		"Version": versionCodeName,
		"UpTime":  uptime.String(),
	})
	photo.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(photo); err != nil {
		log.Ln("Sending message error:", err.Error())
	}

	<-metrika // Send track to Yandex.metrika
}
