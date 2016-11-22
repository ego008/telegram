package main

import (
	"fmt"
	b "github.com/botanio/sdk/go"
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hako/durafmt"
	"github.com/nicksnyder/go-i18n/i18n"
	f "github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	parseMarkdown    = "markdown"
	parseHTML        = "html"
	versionCover     = "AgADAgADO8cxG2OYSwL1_YNSINBpb48ycQ0ABF4-pz5UE6UE1DYCAAEC"
	versionCodeName  = "2.0 \"Busujima Saeko\""
	demonstrationGIF = "BQADAgADNwYAAmOYSwLFYMl_HVAaDwI"
)

var startUptime = time.Now()

func getMessages(message *t.Message) {
	locale := checkLanguage(message.From)

	isCommand := message.IsCommand()
	isPrivate := message.Chat.IsPrivate()
	switch {
	case isCommand /*&& (isPrivate || isMessageToMe)*/ :
		go sendMessages(message, locale)
	case isPrivate && message.From.ID == config.Telegram.Admin && message.Text == "":
		go getTelegramFileID(message) // Admin feature without tracking
	default:
		go easterEggsMessages(message) // Secret actions and commands ;)
	}
}

func sendMessages(message *t.Message, locale i18n.TranslateFunc) {
	lowerCommand := strings.ToLower(message.Command())
	switch lowerCommand {
	case "start": // Requirement Telegram platform
		go startCommand(message, locale)
	case "help": // Requirement Telegram platform
		go helpCommand(message, locale)
	case "settings": // Requirement Telegram platform
		go settingsCommand(message, locale)
	case "cheatsheet":
		go cheatsheetCommand(message, locale)
	case "random":
		go randomCommand(message, locale)
	case "info":
		go infoCommand(message, locale)
	case "donate":
		go donateCommand(message, locale)
	default:
		go eggCommand(message)
	}
}

func startCommand(message *t.Message, locale i18n.TranslateFunc) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/start", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track /start %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	inlineKeyboard := t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonSwitch(locale("button_try"), "hatsune_miku rating:safe"), // Showing tutorial button for demonstration work
		),
	)

	text := locale("message_start", map[string]interface{}{"FirstName": message.From.FirstName})
	reply := t.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func helpCommand(message *t.Message, locale i18n.TranslateFunc) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/help", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track /help %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	document := t.NewDocumentShare(message.Chat.ID, demonstrationGIF)
	if _, err := bot.Send(document); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	text := locale("message_help")
	reply := t.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func settingsCommand(message *t.Message, locale i18n.TranslateFunc) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/settings", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track /settings %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	nsfw := checkNSFW(message.From)

	var nsfwBtn t.InlineKeyboardButton
	if nsfw {
		nsfwBtn = t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(locale("status_on")),
		}), "nsfw_off")
	} else {
		nsfwBtn = t.NewInlineKeyboardButtonData(locale("button_nsfw", map[string]interface{}{
			"Status": strings.ToUpper(locale("status_off")),
		}), "nsfw_on")
	}

	inlineKeyboard := t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			nsfwBtn,
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData(locale("button_language"), "to_lang"),
		),
	)

	text := locale("message_settings")
	reply := t.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func cheatsheetCommand(message *t.Message, locale i18n.TranslateFunc) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/cheatsheet", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track /cheatsheet %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	text := locale("message_cheatsheet")
	reply := t.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func donateCommand(message *t.Message, locale i18n.TranslateFunc) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/donate", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track /donate %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	var donateURL string
	if message.Chat.IsPrivate() {
		donateURL = getBotanURL(message.From.ID, config.Links.Donate)
	} else {
		donateURL = config.Links.Donate
	}

	inlineKeyboard := t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonURL(locale("button_donate"), donateURL),
		),
	)

	text := locale("message_donate")
	reply := t.NewMessage(message.Chat.ID, text)
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func randomCommand(message *t.Message, locale i18n.TranslateFunc) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/random", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track /random %s", answer.Status)
		appMetrika <- true
	})

	nsfw := checkNSFW(message.From)

	if nsfw {
		// Force feedback
		if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatUploadDocument)); err != nil {
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
		bytes := t.FileBytes{
			Name:  randomFile[0].Image,
			Bytes: body,
		}
		uploadFilesProcess(message, bytes, randomFile[0], locale)
	} else {
		// Force feedback
		if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatTyping)); err != nil {
			log.Printf("[Bot] ChatAction send error: %+v", err)
		}

		reply := t.NewMessage(message.Chat.ID, "`¯\\_(ツ)_/¯`")
		reply.ParseMode = parseMarkdown
		reply.DisableWebPagePreview = true
		reply.ReplyToMessageID = message.MessageID
		if _, err := bot.Send(reply); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func infoCommand(message *t.Message, locale i18n.TranslateFunc) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/info", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track /info %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	uptimePeriod := time.Since(startUptime).String()
	uptime, err := durafmt.ParseString(uptimePeriod)
	if err != nil {
		fmt.Printf("[Bot] DuraFmt error: %+v", err)
	}

	inlineKeyboard := t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonURL(locale("button_channel"), config.Links.Channel),
			t.NewInlineKeyboardButtonURL(locale("button_group"), config.Links.Group),
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonURL(locale("button_rate"), config.Links.Rate+bot.Self.UserName),
		),
	)
	photo := t.NewPhotoShare(message.Chat.ID, versionCover)
	photo.Caption = locale("message_info", map[string]interface{}{
		"Version": versionCodeName,
		"UpTime":  uptime.String(),
	})
	photo.ReplyToMessageID = message.MessageID
	photo.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(photo); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}
