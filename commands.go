package main

import (
	"fmt"
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hako/durafmt"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func sendHello(message *tgbotapi.Message) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/start", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /start %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	tutorialButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			// Showing tutorial button for demonstration work
			tgbotapi.NewInlineKeyboardButtonSwitch("See how to do this!", "hatsune_miku rating:safe"),
		),
	)

	reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(startMessage, message.From.FirstName, bot.Self.UserName))
	reply.ParseMode = "markdown"
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &tutorialButton

	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendHelp(message *tgbotapi.Message) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/help", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /help %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, helpMessage)
	reply.ParseMode = "markdown"
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendCheatSheet(message *tgbotapi.Message) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/cheatsheet", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /cheatsheet %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	// For now - get Cheat Sheet from Gelbooru
	reply := tgbotapi.NewMessage(message.Chat.ID, cheatSheetMessage)
	reply.ParseMode = "markdown"
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendRandomPost(message *tgbotapi.Message) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/random", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /random %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadDocument)); err != nil {
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
	_, body, err := fasthttp.Get(nil, randomFile[0].FileURL)
	if err != nil {
		log.Printf("[Bot] Get random image by URL error: %+v", err)
	}
	bytes := tgbotapi.FileBytes{Name: randomFile[0].Image, Bytes: body}

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	if message.Chat.IsPrivate() == true { // Add share-button if chat is private
		originalLink := getBotanURL(message.From.ID, randomFile[0].FileURL)
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Original", originalLink),
				tgbotapi.NewInlineKeyboardButtonSwitch("Share", "id:"+strconv.Itoa(randomFile[0].ID)),
			),
		)
	} else {
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Original image", randomFile[0].FileURL),
			),
		)
	}

	switch {
	case strings.Contains(randomFile[0].FileURL, ".mp4"):
		// Force feedback
		if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadVideo)); err != nil {
			log.Printf("[Bot] ChatAction send error: %+v", err)
		}

		video := tgbotapi.NewVideoUpload(message.Chat.ID, bytes)
		video.ReplyToMessageID = message.MessageID
		video.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(video); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case strings.Contains(randomFile[0].FileURL, ".gif") || strings.Contains(randomFile[0].FileURL, ".webm"):
		// Force feedback
		if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadDocument)); err != nil {
			log.Printf("[Bot] ChatAction send error: %+v", err)
		}

		gif := tgbotapi.NewDocumentUpload(message.Chat.ID, bytes)
		gif.ReplyToMessageID = message.MessageID
		gif.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(gif); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	default:
		// Force feedback
		if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadPhoto)); err != nil {
			log.Printf("[Bot] ChatAction send error: %+v", err)
		}

		image := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
		image.ReplyToMessageID = message.MessageID
		image.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(image); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendBotInfo(message *tgbotapi.Message, startUptime time.Time) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/info", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /info %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	uptimePeriod := time.Since(startUptime).String()
	uptime, err := durafmt.ParseString(uptimePeriod)
	if err != nil {
		fmt.Printf("[Bot] DuraFmt error: %+v", err)
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📢 Channel", config.Telegram.Invite.Channel),
			tgbotapi.NewInlineKeyboardButtonURL("👥 Group", config.Telegram.Invite.Group),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Rate ⭐️⭐️⭐️⭐️⭐️", "https://telegram.me/storebot?start="+bot.Self.UserName),
		),
	)

	photo := tgbotapi.NewPhotoShare(message.Chat.ID, "AgADAgADs8YxG2OYSwJdP213y5L1A68qcQ0ABHvDI3ToOjngT6cBAAEC")
	photo.Caption = fmt.Sprintf(infoMessage, "1.1 Aikawa Jun", uptime.String())
	photo.ReplyToMessageID = message.MessageID
	photo.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(photo); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendDonate(message *tgbotapi.Message) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/donate", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /donate %s", answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	patreonURL := getBotanURL(message.From.ID, "https://patreon.com/toby3d")

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("💸 Become a patron!", patreonURL),
		),
	)

	reply := tgbotapi.NewMessage(message.Chat.ID, donateMessage)
	reply.ParseMode = "markdown"
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendTelegramFileID(message *tgbotapi.Message) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	var fileID string
	switch {
	case message.Audio != nil:
		fileID = message.Audio.FileID
	case message.Document != nil:
		fileID = message.Document.FileID
	case message.Photo != nil:
		photo := *message.Photo
		id := 0
		for i, v := range photo {
			if v.Width > photo[id].Width {
				id = i
			}
		}
		fileID = photo[id].FileID
	case message.Sticker != nil:
		fileID = message.Sticker.FileID
	case message.Video != nil:
		fileID = message.Video.FileID
	case message.Voice != nil:
		fileID = message.Voice.FileID
	}
	if fileID != "" {
		reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("ID: %s", fileID))
		reply.ReplyToMessageID = message.MessageID
		if _, err := bot.Send(reply); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}

// Generate personal tracking-link
func getBotanURL(ID int, URL string) string {
	status, botanURL, err := fasthttp.Get(nil, "https://api.botan.io/s/?token="+config.Botan.Token+"&user_ids="+strconv.Itoa(ID)+"&url="+URL)
	if err != nil || status != 200 {
		log.Printf("[Botan] Generate URL error: %+v", err)
		botanURL = []byte(URL)
	}
	return string(botanURL)
}
