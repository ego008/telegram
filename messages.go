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

const (
	modeMarkdown = "markdown"
	//modeHTML     = "html"
)

func sendSimpleMessage(message *tgbotapi.Message, command string, text string) {
	// Track action
	metrika.TrackAsync(message.From.ID, MetrikaMessage{message}, "/"+command, func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track /%s %s", command, answer.Status)
		appMetrika <- true
	})

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, text)
	reply.ParseMode = modeMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID

	switch command {
	case "start":
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				// Showing tutorial button for demonstration work
				tgbotapi.NewInlineKeyboardButtonSwitch(locale.English.Buttons.FastStart, "hatsune_miku rating:safe"),
			),
		)
		reply.ReplyMarkup = &inlineKeyboard
	case "donate":
		var donateURL string
		if message.Chat.IsPrivate() {
			donateURL = getBotanURL(message.From.ID, config.Links.Donate)
		} else {
			donateURL = config.Links.Donate
		}

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Donate, donateURL),
			),
		)
		reply.ReplyMarkup = &inlineKeyboard
	}

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
	bytes := tgbotapi.FileBytes{
		Name:  randomFile[0].Image,
		Bytes: body,
	}

	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	if message.Chat.IsPrivate() == true { // Add share-button if chat is private
		originalLink := getBotanURL(message.From.ID, randomFile[0].FileURL)
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, originalLink),
				tgbotapi.NewInlineKeyboardButtonSwitch(locale.English.Buttons.Share, "id:"+strconv.Itoa(randomFile[0].ID)),
			),
		)
	} else {
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, randomFile[0].FileURL),
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
			tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Channel, config.Links.Channel),
			tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Group, config.Links.Group),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Rate, config.Links.Rate+bot.Self.UserName),
		),
	)

	photo := tgbotapi.NewPhotoShare(message.Chat.ID, config.Version.Photo)
	photo.Caption = fmt.Sprintf(locale.English.Messages.Info, config.Version.Name, uptime.String())
	photo.ReplyToMessageID = message.MessageID
	photo.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(photo); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendTelegramFileID(message *tgbotapi.Message) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	var uploadFileInfo string
	switch {
	case message.Audio != nil && strings.Contains(message.Audio.MimeType, "ogg") != true: // Upload file As Audio
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Audio.FileID)
	case message.Audio != nil && strings.Contains(message.Audio.MimeType, "ogg") == true: // Upload file as Voice
		if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatRecordAudio)); err != nil {
			log.Printf("[Bot] ChatAction send error: %+v", err)
		}

		url, err := bot.GetFileDirectURL(message.Audio.FileID)
		if err != nil {
			log.Printf("ERROR: %+v", err)
		}

		_, body, err := fasthttp.Get(nil, url)
		if err != nil {
			log.Printf("Get file error: %+v", err)
		}
		bytes := tgbotapi.FileBytes{Name: message.Audio.FileID, Bytes: body}

		voice := tgbotapi.NewVoiceUpload(message.Chat.ID, bytes)
		voice.Duration = message.Audio.Duration
		voice.ReplyToMessageID = message.MessageID
		resp, err := bot.Send(voice)
		if err != nil {
			log.Printf("Sending message error: %+v", err)
		}

		uploadFileInfo = fmt.Sprintf("ID: %s\nDuration: %d", resp.Voice.FileID, resp.Voice.Duration)
	case message.Document != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Document.FileID)
	case message.Photo != nil: // Get large file ID
		photo := *message.Photo
		id := 0
		for i, v := range photo {
			if v.Width > photo[id].Width {
				id = i
			}
		}
		uploadFileInfo = fmt.Sprintf("ID: %s", photo[id].FileID)
	case message.Sticker != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Sticker.FileID)
	case message.Video != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Video.FileID)
	case message.Voice != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Voice.FileID)
	}
	if uploadFileInfo != "" {
		reply := tgbotapi.NewMessage(message.Chat.ID, uploadFileInfo)
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
