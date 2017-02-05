package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
	f "github.com/valyala/fasthttp"
	tg "gopkg.in/telegram-bot-api.v4"
)

func uploadFilesProcess(message *tg.Message, bytes tg.FileBytes, randomFile Post, locale i18n.TranslateFunc) {
	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatUploadDocument)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	var inlineKeyboard tg.InlineKeyboardMarkup
	if message.Chat.IsPrivate() { // Add share-button if chat is private
		originalLink := getBotanURL(message.From.ID, randomFile.FileURL)
		inlineKeyboard = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonURL(locale("button_original"), originalLink),
				tg.NewInlineKeyboardButtonSwitch(locale("button_share"), "id:"+strconv.Itoa(randomFile.ID)),
			),
		)
	} else {
		inlineKeyboard = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonURL(locale("button_original"), randomFile.FileURL),
			),
		)
	}

	switch {
	case strings.Contains(randomFile.FileURL, ".mp4"):
		video := tg.NewVideoUpload(message.Chat.ID, bytes)
		video.ReplyToMessageID = message.MessageID
		video.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(video); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case strings.Contains(randomFile.FileURL, ".gif"):
		gif := tg.NewDocumentUpload(message.Chat.ID, bytes)
		gif.ReplyToMessageID = message.MessageID
		gif.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(gif); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case strings.Contains(randomFile.FileURL, ".webm"):
		pageURL := BlushBoard + "/hash/" + randomFile.Hash
		text := locale("message_blushboard", map[string]interface{}{
			"Type":  strings.Title(locale("type_video")),
			"Owner": randomFile.Owner,
			"URL":   pageURL,
		})
		reply := tg.NewMessage(message.Chat.ID, text)
		reply.ParseMode = parseMarkdown
		reply.DisableWebPagePreview = false
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(reply); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	default:
		image := tg.NewPhotoUpload(message.Chat.ID, bytes)
		image.ReplyToMessageID = message.MessageID
		image.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(image); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}

func getTelegramFileID(message *tg.Message) {
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
	reply := tg.NewMessage(message.Chat.ID, uploadFileInfo)
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

func getLargePhoto(message *tg.Message) string {
	photo := *message.Photo
	id := 0
	for i, v := range photo {
		if v.Width > photo[id].Width {
			id = i
		}
	}
	return photo[id].FileID
}

func getVoiceFromAudio(message *tg.Message) string {
	if _, err := bot.Send(tg.NewChatAction(message.Chat.ID, tg.ChatRecordAudio)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	url, err := bot.GetFileDirectURL(message.Audio.FileID)
	if err != nil {
		log.Printf("ERROR: %+v", err)
	}

	_, body, err := f.Get(nil, url)
	if err != nil {
		log.Printf("Get file error: %+v", err)
	}
	bytes := tg.FileBytes{
		Name:  message.Audio.FileID,
		Bytes: body,
	}

	voice := tg.NewVoiceUpload(message.Chat.ID, bytes)
	voice.Duration = message.Audio.Duration
	voice.ReplyToMessageID = message.MessageID
	resp, err := bot.Send(voice)
	if err != nil {
		log.Printf("Sending message error: %+v", err)
	}

	return resp.Voice.FileID
}

// Generate personal tracking-link
func getBotanURL(id int, url string) string {
	req := fmt.Sprintf("https://api.botan.io/s/?token=%s&user_ids=%d&url=%s", cfg["botan"].(string), id, url)
	status, botanURL, err := f.Get(nil, req)
	if err != nil || status != 200 {
		log.Printf("[Botan] Generate URL error: %+v", err)
		botanURL = []byte(url)
	}
	return string(botanURL)
}
