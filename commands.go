package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func sendHello(message *tgbotapi.Message) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error^ %+v", err)
	}

	answer := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(startMessage, message.From.FirstName))
	answer.ParseMode = "markdown"
	answer.DisableWebPagePreview = true
	answer.ReplyToMessageID = message.MessageID
	if message.Chat.IsPrivate() == true {
		answer.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				// Showing tutorial button only in private chat for demonstration work
				tgbotapi.NewInlineKeyboardButtonSwitch("See how to do this!", "hatsune_miku rating:safe"),
			),
		)
	}
	if _, err := bot.Send(answer); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

func sendHelp(message *tgbotapi.Message) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error^ %+v", err)
	}

	answer := tgbotapi.NewMessage(message.Chat.ID, helpMessage)
	answer.ParseMode = "markdown"
	answer.DisableWebPagePreview = true
	answer.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(answer); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

func sendCheatSheet(message *tgbotapi.Message) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error^ %+v", err)
	}

	// For now - get Cheat Sheet from Gelbooru
	answer := tgbotapi.NewMessage(message.Chat.ID, cheatSheetMessage)
	answer.ParseMode = "markdown"
	answer.DisableWebPagePreview = true
	answer.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(answer); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

func sendRandomPost(message *tgbotapi.Message) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadDocument)); err != nil {
		log.Printf("[Bot] ChatAction send error^ %+v", err)
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
		log.Println("[Bot] This is not image. Reroll dice!")
	}

	var button tgbotapi.InlineKeyboardMarkup
	if message.Chat.IsPrivate() == true { // Add share-button if chat is private
		button = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Original image", randomFile[0].FileURL),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonSwitch("Share", "id:"+strconv.Itoa(randomFile[0].ID)),
			),
		)
	} else {
		button = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Original image", randomFile[0].FileURL),
			),
		)
	}

	_, body, err := fasthttp.Get(nil, randomFile[0].FileURL)
	if err != nil {
		log.Printf("[Bot] Get random image by URL error: %+v", err)
	}
	bytes := tgbotapi.FileBytes{Name: randomFile[0].Image, Bytes: body}

	switch {
	case strings.Contains(randomFile[0].FileURL, ".mp4"):
		// Force feedback
		if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadVideo)); err != nil {
			log.Printf("[Bot] ChatAction send error^ %+v", err)
		}

		video := tgbotapi.NewVideoUpload(message.Chat.ID, bytes)
		video.ReplyToMessageID = message.MessageID
		video.ReplyMarkup = &button
		if _, err := bot.Send(video); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	case strings.Contains(randomFile[0].FileURL, ".gif") || strings.Contains(randomFile[0].FileURL, ".webm"):
		// Force feedback
		if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadDocument)); err != nil {
			log.Printf("[Bot] ChatAction send error^ %+v", err)
		}

		gif := tgbotapi.NewDocumentUpload(message.Chat.ID, bytes)
		gif.ReplyToMessageID = message.MessageID
		gif.ReplyMarkup = &button
		if _, err := bot.Send(gif); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	default:
		// Force feedback
		if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatUploadPhoto)); err != nil {
			log.Printf("[Bot] ChatAction send error^ %+v", err)
		}

		image := tgbotapi.NewPhotoUpload(message.Chat.ID, bytes)
		image.ReplyToMessageID = message.MessageID
		image.ReplyMarkup = &button
		if _, err := bot.Send(image); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}

func sendBotInfo(message *tgbotapi.Message, startUptime time.Time) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error^ %+v", err)
	}

	uptime := time.Since(startUptime).String()

	// For now - get Cheat Sheet from Gelbooru
	answer := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(infoMessage, uptime))
	answer.ParseMode = "markdown"
	answer.DisableWebPagePreview = true
	answer.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(answer); err != nil {
		log.Printf("[Bot] Sending message error: %+v", err)
	}
}

func sendTelegramFileID(message *tgbotapi.Message) {
	// Force feedback
	if _, err := bot.Send(tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)); err != nil {
		log.Printf("[Bot] ChatAction send error^ %+v", err)
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
		answer := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("ID: %s", fileID))
		answer.ReplyToMessageID = message.MessageID
		if _, err := bot.Send(answer); err != nil {
			log.Printf("[Bot] Sending message error: %+v", err)
		}
	}
}
