package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hako/durafmt"
	"github.com/nicksnyder/go-i18n/i18n"
	f "github.com/valyala/fasthttp"
)

const (
	parseMarkdown    = "markdown"
	parseHTML        = "html"
	versionCover     = "AgADAgADO8cxG2OYSwL1_YNSINBpb48ycQ0ABF4-pz5UE6UE1DYCAAEC"
	versionCodeName  = "2.0 \"Busujima Saeko\""
	demonstrationGIF = "BQADAgADNwYAAmOYSwLFYMl_HVAaDwI"
)

var startUptime = time.Now()

func getMessages(message *tg.Message) {
	isCommand := message.IsCommand()
	isPrivate := message.Chat.IsPrivate()
	switch {
	case isCommand /*&& (isPrivate || isMessageToMe)*/ :
		go sendMessages(message, T)
	case isPrivate && message.From.ID == int(cfg["telegram_admin"].(float64)) && message.Text == "":
		go getTelegramFileID(message) // Admin feature without tracking
	default:
		go easterEggsMessages(message) // Secret actions and commands ;)
	}
}

func sendMessages(message *tg.Message, T i18n.TranslateFunc) {
	lowerCommand := strings.ToLower(message.Command())
	switch lowerCommand {
	case "start": // Requirement Telegram platform
		go startCommand(message, T)
	case "help": // Requirement Telegram platform
		go helpCommand(message, T)
	// case "settings": // Requirement Telegram platform
	// 	go settingsCommand(message, T)
	case "cheatsheet":
		go cheatsheetCommand(message, T)
	case "random":
		go randomCommand(message, T)
	case "info":
		go infoCommand(message, T)
	case "donate":
		go donateCommand(message, T)
	default:
		go eggCommand(message)
	}
}

func startCommand(message *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(message.From.ID, MetrikaMessage{message}, "/start", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /start %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	inlineKeyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonSwitch(T("button_try"), "hatsune_miku rating:safe"), // Showing tutorial button for demonstration work
		),
	)

	text := T("message_start", map[string]interface{}{"FirstName": message.From.FirstName})
	reply := tg.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-metrika // Send track to Yandex.metrika
}

func helpCommand(message *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(message.From.ID, MetrikaMessage{message}, "/help", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /help %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	document := tg.NewDocumentShare(message.Chat.ID, demonstrationGIF)
	if _, err := bot.Send(document); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	text := T("message_help")
	reply := tg.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-metrika // Send track to Yandex.metrika
}

/*
func settingsCommand(message *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(message.From.ID, MetrikaMessage{message}, "/settings", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /settings %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	var nsfwBtn tg.InlineKeyboardButton
	if nsfw {
		nsfwBtn = tg.NewInlineKeyboardButtonData(T("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(T("status_on")),
		}), "nsfw_off")
	} else {
		nsfwBtn = tg.NewInlineKeyboardButtonData(T("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(T("status_off")),
		}), "nsfw_on")
	}

	inlineKeyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			nsfwBtn,
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(T("button_language"), "to_lang"),
		),
	)

	text := T("message_settings")
	reply := tg.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-metrika // Send track to Yandex.metrika
}
*/

func cheatsheetCommand(message *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(message.From.ID, MetrikaMessage{message}, "/cheatsheet", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /cheatsheet %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	text := T("message_cheatsheet")
	reply := tg.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-metrika // Send track to Yandex.metrika
}

func donateCommand(message *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(message.From.ID, MetrikaMessage{message}, "/donate", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /donate %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	var donateURL string
	if message.Chat.IsPrivate() {
		donateURL = getBotanURL(message.From.ID, cfg["link_donate"].(string))
	} else {
		donateURL = cfg["link_donate"].(string)
	}

	inlineKeyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(T("button_donate"), donateURL),
		),
	)

	text := T("message_donate")
	reply := tg.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-metrika // Send track to Yandex.metrika
}

func randomCommand(message *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(message.From.ID, MetrikaMessage{message}, "/random", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /random %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatUploadDocument)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
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
		log.Println("[Bot] Nothing. Reroll dice!")
	}

	log.Printf("[Bot] Get random file by URL: %s", randomFile[0].FileURL)
	_, body, err := f.Get(nil, randomFile[0].FileURL)
	if err != nil {
		log.Printf("[Bot] Get random image by URL error: %+v", err)
	}
	bytes := tg.FileBytes{
		Name:  randomFile[0].Image,
		Bytes: body,
	}
	uploadFilesProcess(message, bytes, randomFile[0], T)

	<-metrika // Send track to Yandex.metrika
}

func infoCommand(message *tg.Message, T i18n.TranslateFunc) {
	// Track action
	b.TrackAsync(message.From.ID, MetrikaMessage{message}, "/info", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /info %s", answer.Status)
		metrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	uptimePeriod := time.Since(startUptime).String()
	uptime, err := durafmt.ParseString(uptimePeriod)
	if err != nil {
		fmt.Printf("[Bot] DuraFmt error: %+v", err)
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
	photo := tg.NewPhotoShare(message.Chat.ID, versionCover)
	photo.Caption = T("message_info", map[string]interface{}{
		"Version": versionCodeName,
		"UpTime":  uptime.String(),
	})
	photo.ReplyToMessageID = message.MessageID
	photo.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(photo); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-metrika // Send track to Yandex.metrika
}
