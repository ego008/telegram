package main

import (
	"fmt"
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

var results []interface{}

const BlushBoard = "http://beta.hentaidb.pw"

func getInlineResults(cacheTime int, inline *tgbotapi.InlineQuery) {
	// Track action
	metrika.TrackAsync(inline.From.ID, MetrikaInlineQuery{inline}, "Search", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track Search %s", answer.Status)
		appMetrika <- true
	})

	// Check result pages
	var post []Post
	var resultPage int
	if len(inline.Offset) > 0 {
		resultPage, _ = strconv.Atoi(inline.Offset)
		post = getPosts(
			Request{
				Limit:  50,
				PageID: resultPage,
				Tags:   inline.Query,
			})
	} else {
		post = getPosts(
			Request{
				Limit: 50,
				Tags:  inline.Query,
			})
	}

	results = nil
	switch {
	case len(post) > 0:
		for i := 0; i < len(post); i++ {
			inlineResult(post[i])
		}
	case len(post) == 0: // Found nothing
		emptyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Channel, config.Links.Channel),
				tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Group, config.Links.Group),
			),
		)
		emptyMessage := fmt.Sprintf("*%s*\n%s", locale.English.Inline.NoResult.Title, locale.English.Inline.NoResult.Description)
		empty := tgbotapi.NewInlineQueryResultArticleMarkdown(inline.ID, locale.English.Inline.NoResult.Title, emptyMessage)
		empty.Description = locale.English.Inline.NoResult.Description
		empty.ReplyMarkup = &emptyKeyboard
		results = append(results, empty)
	}

	// Configure inline-mode
	inlineConfig := tgbotapi.InlineConfig{}
	inlineConfig.InlineQueryID = inline.ID
	inlineConfig.IsPersonal = true
	inlineConfig.CacheTime = cacheTime
	inlineConfig.Results = results
	// If available next page of results
	if len(post) == 50 {
		resultPage++
		inlineConfig.NextOffset = strconv.Itoa(resultPage)
	}

	if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
		log.Printf("[Bot] Answer inline-query error: %+v", err)
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func inlineResult(post Post) {
	// Universal(?) preview url
	preview := config.Resource[20].Settings.URL + config.Resource[20].Settings.ThumbsDir + post.Directory + config.Resource[20].Settings.ThumbsPart + post.Hash + ".jpg"
	post.Rating = setResultRating(post.Rating)
	resultKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, post.FileURL),
		),
	)

	switch {
	case strings.Contains(post.FileURL, ".webm"): // Not support yet. Show tip about "hidden" result
		video := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(post.ID), BlushBoard+"/embed/"+post.Hash)
		video.MimeType = "text/html"
		video.ThumbURL = preview
		video.Title = fmt.Sprintf(locale.English.Inline.Result.Title, strings.Title(locale.English.Types.Video), post.Owner)
		video.Width = post.Width
		video.Height = post.Height
		video.Description = fmt.Sprintf(locale.English.Inline.Result.Description, post.Rating, post.Tags)
		videoURL := BlushBoard + "/hash/" + post.Hash
		video.InputMessageContent = tgbotapi.InputTextMessageContent{
			Text:                  fmt.Sprintf(locale.English.Messages.BlushBoard, strings.Title(locale.English.Types.Video), post.Owner, videoURL),
			ParseMode:             parseMarkdown,
			DisableWebPagePreview: false,
		}
		results = append(results, video)
	case strings.Contains(post.FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
		video := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(post.ID), post.FileURL)
		video.MimeType = "video/mp4"
		video.ThumbURL = preview
		video.Title = fmt.Sprintf(locale.English.Inline.Result.Title, strings.Title(locale.English.Types.Video), post.Owner)
		video.Width = post.Width
		video.Height = post.Height
		video.Description = fmt.Sprintf(locale.English.Inline.Result.Description, post.Rating, post.Tags)
		video.ReplyMarkup = &resultKeyboard
		results = append(results, video)
	case strings.Contains(post.FileURL, ".gif"):
		gif := tgbotapi.NewInlineQueryResultGIF(strconv.Itoa(post.ID), post.FileURL)
		gif.Width = post.Width
		gif.Height = post.Height
		gif.ThumbURL = post.FileURL
		gif.Title = fmt.Sprintf(locale.English.Inline.Result.Title, strings.Title(locale.English.Types.Animation), post.Owner)
		gif.ReplyMarkup = &resultKeyboard
		results = append(results, gif)
	default:
		image := tgbotapi.NewInlineQueryResultPhoto(strconv.Itoa(post.ID), post.FileURL)
		image.ThumbURL = preview
		image.Width = post.Width
		image.Height = post.Height
		image.Title = fmt.Sprintf(locale.English.Inline.Result.Title, strings.Title(locale.English.Types.Image), post.Owner)
		image.Description = fmt.Sprintf(locale.English.Inline.Result.Description, post.Rating, post.Tags)
		image.ReplyMarkup = &resultKeyboard
		results = append(results, image)
	}
}

func setResultRating(rating string) string {
	switch rating {
	case "s":
		return locale.English.Rating.Safe
	case "e":
		return locale.English.Rating.Explicit
	case "q":
		return locale.English.Rating.Questionable
	default:
		return locale.English.Rating.Unknown
	}
}

func sendInlineResult(result *tgbotapi.ChosenInlineResult) {
	metrika.TrackAsync(result.From.ID, MetrikaChosenInlineResult{result}, "Find", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track Find %s", answer.Status)
		appMetrika <- true
	})

	<-appMetrika // Send track to Yandex.AppMetrika
}
