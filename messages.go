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
	parseMarkdown    = "markdown"
	parseHTML        = "html"
	versionCover     = "AgADAgADO8cxG2OYSwL1_YNSINBpb48ycQ0ABF4-pz5UE6UE1DYCAAEC"
	versionCodeName  = "2.0 \"Busujima Saeko\""
	demonstrationGIF = "BQADAgADNwYAAmOYSwLFYMl_HVAaDwI"
)

var startUptime = time.Now()

func sendMessages(message *tgbotapi.Message) {
	isCommand := message.IsCommand()
	isPrivate := message.Chat.IsPrivate()
	isMessageToMe := bot.IsMessageToMe(*message)
	switch {
	case isCommand && (isPrivate || isMessageToMe):
		go commandActions(message)
	case isPrivate && message.From.ID == config.Telegram.Admin && message.Text == "":
		go getTelegramFileID(message) // Admin feature without tracking
	default:
		go getEggMessage(message) // Secret actions and commands ;)
	}
}

func commandActions(message *tgbotapi.Message) {
	lowerCommand := strings.ToLower(message.Command())
	switch lowerCommand {
	case "start": // Requirement Telegram platform
		messageText := fmt.Sprintf(locale.English.Messages.Start, message.From.FirstName, bot.Self.UserName)
		go sendSimpleMessage(message, lowerCommand, messageText)
	case "help": // Requirement Telegram platform
		go sendSimpleMessage(message, lowerCommand, locale.English.Messages.Help)
	case "cheatsheet":
		go sendSimpleMessage(message, lowerCommand, locale.English.Messages.CheatSheet)
	case "random":
		go sendRandomPost(message)
	case "info":
		go sendBotInfo(message)
	case "donate":
		go sendSimpleMessage(message, lowerCommand, locale.English.Messages.Donate)
	default:
		go sendEggAction(message)
	}
}

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
	reply.ParseMode = parseMarkdown
	reply.DisableWebPagePreview = true
	reply.ReplyToMessageID = message.MessageID

	switch command {
	case "start":
		document := tgbotapi.NewDocumentShare(message.Chat.ID, demonstrationGIF)
		if _, err := bot.Send(document); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}

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

	if command == "help" {
		document := tgbotapi.NewDocumentShare(message.Chat.ID, demonstrationGIF)
		if _, err := bot.Send(document); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
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
	uploadFilesProcess(message, bytes, randomFile[0])

	<-appMetrika // Send track to Yandex.AppMetrika
}

func sendBotInfo(message *tgbotapi.Message) {
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
	photo := tgbotapi.NewPhotoShare(message.Chat.ID, versionCover)
	photo.Caption = fmt.Sprintf(locale.English.Messages.Info, versionCodeName, uptime.String())
	photo.ReplyToMessageID = message.MessageID
	photo.ReplyMarkup = &inlineKeyboard
	if _, err := bot.Send(photo); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func uploadFilesProcess(message *tgbotapi.Message, bytes tgbotapi.FileBytes, randomFile Post) {
	var inlineKeyboard tgbotapi.InlineKeyboardMarkup
	if message.Chat.IsPrivate() == true { // Add share-button if chat is private
		originalLink := getBotanURL(message.From.ID, randomFile.FileURL)
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, originalLink),
				tgbotapi.NewInlineKeyboardButtonSwitch(locale.English.Buttons.Share, "id:"+strconv.Itoa(randomFile.ID)),
			),
		)
	} else {
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, randomFile.FileURL),
			),
		)
	}

	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadDocument)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	switch {
	case strings.Contains(randomFile.FileURL, ".mp4"):
		video := tgbotapi.NewVideoUpload(message.Chat.ID, bytes)
		video.ReplyToMessageID = message.MessageID
		video.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(video); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case strings.Contains(randomFile.FileURL, ".gif") || strings.Contains(randomFile.FileURL, ".webm"):
		gif := tgbotapi.NewDocumentUpload(message.Chat.ID, bytes)
		gif.ReplyToMessageID = message.MessageID
		gif.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(gif); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	default:
		image := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
		image.ReplyToMessageID = message.MessageID
		image.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(image); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}

func getTelegramFileID(message *tgbotapi.Message) {
	var uploadFileInfo string
	switch {
	case message.Audio != nil: // Upload file as Voice
		if strings.Contains(message.Audio.MimeType, "ogg") == true {
			voice := getVoiceFromAudio(message)
			uploadFileInfo = fmt.Sprintf("ID: %s", voice)
		} else {
			uploadFileInfo = fmt.Sprintf("ID: %s", message.Audio.FileID)
		}
	case message.Document != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Document.FileID)
	case message.Photo != nil: // Get large file ID
		photo := getLargePhoto(message)
		uploadFileInfo = fmt.Sprintf("ID: %s", photo)
	case message.Sticker != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Sticker.FileID)
	case message.Video != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Video.FileID)
	case message.Voice != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", message.Voice.FileID)
	}
	reply := tgbotapi.NewMessage(message.Chat.ID, uploadFileInfo)
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

func getLargePhoto(message *tgbotapi.Message) string {
	photo := *message.Photo
	id := 0
	for i, v := range photo {
		if v.Width > photo[id].Width {
			id = i
		}
	}
	return photo[id].FileID
}

func getVoiceFromAudio(message *tgbotapi.Message) string {
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
	bytes := tgbotapi.FileBytes{
		Name:  message.Audio.FileID,
		Bytes: body,
	}

	voice := tgbotapi.NewVoiceUpload(message.Chat.ID, bytes)
	voice.Duration = message.Audio.Duration
	voice.ReplyToMessageID = message.MessageID
	resp, err := bot.Send(voice)
	if err != nil {
		log.Printf("Sending message error: %+v", err)
	}

	return resp.Voice.FileID
}

// Generate personal tracking-link
func getBotanURL(ID int, URL string) string {
	const botan = "https://api.botan.io/s/"
	status, botanURL, err := fasthttp.Get(nil, botan+"?token="+config.Botan.Token+"&user_ids="+strconv.Itoa(ID)+"&url="+URL)
	if err != nil || status != 200 {
		log.Printf("[Botan] Generate URL error: %+v", err)
		botanURL = []byte(URL)
	}
	return string(botanURL)
}
