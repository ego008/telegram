package main

import (
	"fmt"
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	f "github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
)

func uploadFilesProcess(message *t.Message, bytes t.FileBytes, randomFile Post) {
	// lang := checkLanguage(message.From)

	// Force feedback
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatUploadDocument)); err != nil {
		log.Printf("[Bot] ChatAction send error: %+v", err)
	}

	var inlineKeyboard t.InlineKeyboardMarkup
	if message.Chat.IsPrivate() { // Add share-button if chat is private
		originalLink := getBotanURL(message.From.ID, randomFile.FileURL)
		inlineKeyboard = t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, originalLink),
				t.NewInlineKeyboardButtonSwitch(locale.English.Buttons.Share, "id:"+strconv.Itoa(randomFile.ID)),
			),
		)
	} else {
		inlineKeyboard = t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, randomFile.FileURL),
			),
		)
	}

	switch {
	case strings.Contains(randomFile.FileURL, ".mp4"):
		video := t.NewVideoUpload(message.Chat.ID, bytes)
		video.ReplyToMessageID = message.MessageID
		video.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(video); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case strings.Contains(randomFile.FileURL, ".gif"):
		gif := t.NewDocumentUpload(message.Chat.ID, bytes)
		gif.ReplyToMessageID = message.MessageID
		gif.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(gif); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case strings.Contains(randomFile.FileURL, ".webm"):
		pageURL := BlushBoard + "/hash/" + randomFile.Hash
		text := fmt.Sprintf(locale.English.Messages.BlushBoard, strings.Title(locale.English.Types.Video), randomFile.Owner, pageURL)
		reply := t.NewMessage(message.Chat.ID, text)
		reply.ParseMode = parseMarkdown
		reply.DisableWebPagePreview = false
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(reply); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	default:
		image := t.NewPhotoUpload(message.Chat.ID, bytes)
		image.ReplyToMessageID = message.MessageID
		image.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(image); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}

func getTelegramFileID(message *t.Message) {
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
	reply := t.NewMessage(message.Chat.ID, uploadFileInfo)
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

func getLargePhoto(message *t.Message) string {
	photo := *message.Photo
	id := 0
	for i, v := range photo {
		if v.Width > photo[id].Width {
			id = i
		}
	}
	return photo[id].FileID
}

func getVoiceFromAudio(message *t.Message) string {
	if _, err := bot.Send(t.NewChatAction(message.Chat.ID, t.ChatRecordAudio)); err != nil {
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
	bytes := t.FileBytes{
		Name:  message.Audio.FileID,
		Bytes: body,
	}

	voice := t.NewVoiceUpload(message.Chat.ID, bytes)
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
	req := fmt.Sprintf("https://api.botan.io/s/?token=%s&user_ids=%d&url=%s", config.Botan.Token, id, url)
	status, botanURL, err := f.Get(nil, req)
	if err != nil || status != 200 {
		log.Printf("[Botan] Generate URL error: %+v", err)
		botanURL = []byte(url)
	}
	return string(botanURL)
}
