package main

import (
	"fmt"
	"strconv"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
	f "github.com/valyala/fasthttp"
)

func uploadFilesProcess(msg *tg.Message, bytes tg.FileBytes, randomFile Post, T i18n.TranslateFunc) {
	// Force feedback
	if _, err := bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatUploadDocument)); err != nil {
		log.Ln("ChatAction send error:", err.Error())
	}

	var inlineKeyboard tg.InlineKeyboardMarkup
	if msg.Chat.IsPrivate() { // Add share-button if chat is private
		originalLink := getBotanURL(msg.From.ID, fmt.Sprint("https:", randomFile.FileURL))
		inlineKeyboard = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonURL(T("button_original"), originalLink),
				tg.NewInlineKeyboardButtonSwitch(T("button_share"), "id:"+strconv.Itoa(randomFile.ID)),
			),
		)
	} else {
		inlineKeyboard = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonURL(T("button_original"), fmt.Sprint("https:", randomFile.FileURL)),
			),
		)
	}

	switch {
	case strings.Contains(fmt.Sprint("https:", randomFile.FileURL), ".mp4"):
		video := tg.NewVideoUpload(msg.Chat.ID, bytes)
		video.ReplyToMessageID = msg.MessageID
		video.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(video); err != nil {
			log.Ln("[Bot] Sending message error:", err.Error())
		}
	case strings.Contains(fmt.Sprint("https:", randomFile.FileURL), ".gif"):
		gif := tg.NewDocumentUpload(msg.Chat.ID, bytes)
		gif.ReplyToMessageID = msg.MessageID
		gif.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(gif); err != nil {
			log.Ln("[Bot] Sending message error:", err.Error())
		}
	case strings.Contains(fmt.Sprint("https:", randomFile.FileURL), ".webm"):
		pageURL := BlushBoard + "/hash/" + randomFile.Hash
		text := T("message_blushboard", map[string]interface{}{
			"Type":  strings.Title(T("type_video")),
			"Owner": randomFile.Owner,
			"URL":   pageURL,
		})
		reply := tg.NewMessage(msg.Chat.ID, text)
		reply.ParseMode = parseMarkdown
		reply.DisableWebPagePreview = false
		reply.ReplyToMessageID = msg.MessageID
		reply.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(reply); err != nil {
			log.Ln("[Bot] Sending message error:", err.Error())
		}
	default:
		image := tg.NewPhotoUpload(msg.Chat.ID, bytes)
		image.ReplyToMessageID = msg.MessageID
		image.ReplyMarkup = &inlineKeyboard
		if _, err := bot.Send(image); err != nil {
			log.Ln("[Bot] Sending message error:", err.Error())
		}
	}
}

func getTelegramFileID(msg *tg.Message) {
	var uploadFileInfo string
	switch {
	case msg.Audio != nil: // Upload file as Voice
		if strings.Contains(msg.Audio.MimeType, "ogg") == true {
			voice := getVoiceFromAudio(msg)
			uploadFileInfo = fmt.Sprintf("ID: %s", voice)
		} else {
			uploadFileInfo = fmt.Sprintf("ID: %s", msg.Audio.FileID)
		}
	case msg.Document != nil:
		uploadFileInfo = fmt.Sprintf("ID: %s", msg.Document.FileID)
	case msg.Photo != nil: // Get large file ID
		photo := getLargePhoto(msg)
		uploadFileInfo = fmt.Sprintf("ID: %s", photo)
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
		log.Ln("[Bot] Sending message error:", err.Error())
	}
}

func getLargePhoto(msg *tg.Message) string {
	photo := *msg.Photo
	id := 0
	for i, v := range photo {
		if v.Width > photo[id].Width {
			id = i
		}
	}
	return photo[id].FileID
}

func getVoiceFromAudio(msg *tg.Message) string {
	if _, err := bot.Send(tg.NewChatAction(msg.Chat.ID, tg.ChatRecordAudio)); err != nil {
		log.Ln("[Bot] ChatAction send error:", err.Error())
	}

	url, err := bot.GetFileDirectURL(msg.Audio.FileID)
	if err != nil {
		log.Ln("ERROR:", err.Error())
	}

	_, body, err := f.Get(nil, url)
	if err != nil {
		log.Ln("Get file error:", err.Error())
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
		log.Ln("Sending message error:", err.Error())
	}

	return resp.Voice.FileID
}

// Generate personal tracking-link
func getBotanURL(id int, url string) string {
	req := fmt.Sprintf("https://api.botan.io/s/?token=%s&user_ids=%d&url=%s", cfg["botan"].(string), id, url)
	status, botanURL, err := f.Get(nil, req)
	if err != nil || status != 200 {
		log.Ln("[Botan] Generate URL error:", err.Error())
		botanURL = []byte(url)
	}
	return string(botanURL)
}
