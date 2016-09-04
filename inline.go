package main

import (
	"fmt"
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

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
		post = getPosts(Request{Limit: 50, PageID: resultPage, Tags: inline.Query})
	} else {
		post = getPosts(Request{Limit: 50, Tags: inline.Query})
	}

	// Analysis of results
	var result []interface{}
	switch {
	case len(post) > 0:
		for i := 0; i < len(post); i++ {
			// Universal(?) preview url
			preview := config.Resource[20].Settings.URL + config.Resource[20].Settings.ThumbsDir + post[i].Directory + config.Resource[20].Settings.ThumbsPart + post[i].Hash + ".jpg"

			// Rating
			var rating string
			switch post[i].Rating {
			case "s":
				rating = locale.English.Rating.Safe
			case "e":
				rating = locale.English.Rating.Explicit
			case "q":
				rating = locale.English.Rating.Questionable
			default:
				rating = locale.English.Rating.Unknown
			}

			resultKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, post[i].FileURL),
				),
			)

			switch {
			case strings.Contains(post[i].FileURL, ".webm"): // Not support yet. Show tip about "hidden" result
				videoKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Original, post[i].FileURL),
						tgbotapi.NewInlineKeyboardButtonURL(locale.English.Buttons.Donate, config.Links.Donate),
					),
				)

				messageText := fmt.Sprintf("*%s*\n%s", locale.English.Inline.HiddenResult.Title, locale.English.Inline.HiddenResult.Description)
				video := tgbotapi.NewInlineQueryResultArticleMarkdown(strconv.Itoa(post[i].ID), locale.English.Inline.HiddenResult.Title, messageText)
				video.URL = post[i].FileURL
				video.Description = locale.English.Inline.HiddenResult.Description
				video.ThumbURL = preview
				video.ThumbWidth = post[i].Width
				video.ThumbHeight = post[i].Height
				video.ReplyMarkup = &videoKeyboard
				result = append(result, video)
			case strings.Contains(post[i].FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
				video := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(post[i].ID), post[i].FileURL)
				video.MimeType = "video/mp4"
				video.ThumbURL = preview
				video.Width = post[i].Width
				video.Height = post[i].Height
				video.Title = fmt.Sprintf(locale.English.Inline.Result.Title, strings.Title(locale.English.Types.Video), post[i].Owner)
				video.Description = fmt.Sprintf(locale.English.Inline.Result.Description, &rating, post[i].Tags)
				video.ReplyMarkup = &resultKeyboard
				result = append(result, video)
			case strings.Contains(post[i].FileURL, ".gif"):
				gif := tgbotapi.NewInlineQueryResultGIF(strconv.Itoa(post[i].ID), post[i].FileURL)
				gif.ThumbURL = post[i].FileURL
				gif.Width = post[i].Width
				gif.Height = post[i].Height
				gif.Title = fmt.Sprintf(locale.English.Inline.Result.Title, strings.Title(locale.English.Types.Animation), post[i].Owner)
				gif.ReplyMarkup = &resultKeyboard
				result = append(result, gif)
			default:
				image := tgbotapi.NewInlineQueryResultPhoto(strconv.Itoa(post[i].ID), post[i].FileURL)
				image.ThumbURL = preview
				image.Width = post[i].Width
				image.Height = post[i].Height
				image.Title = fmt.Sprintf(locale.English.Inline.Result.Title, strings.Title(locale.English.Types.Image), post[i].Owner)
				image.Description = fmt.Sprintf(locale.English.Inline.Result.Description, &rating, post[i].Tags)
				image.ReplyMarkup = &resultKeyboard
				result = append(result, image)
			}
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
		result = append(result, empty)
	}

	// Configure inline-mode
	inlineConfig := tgbotapi.InlineConfig{}
	inlineConfig.InlineQueryID = inline.ID
	inlineConfig.IsPersonal = true
	inlineConfig.CacheTime = cacheTime
	inlineConfig.Results = result
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

func sendInlineResult(result *tgbotapi.ChosenInlineResult) {
	metrika.TrackAsync(result.From.ID, MetrikaChosenInlineResult{result}, "Find", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track Find %s", answer.Status)
		appMetrika <- true
	})

	<-appMetrika // Send track to Yandex.AppMetrika
}

func getCallbackAction(callback *tgbotapi.CallbackQuery) {
	metrika.TrackAsync(callback.From.ID, MetrikaCallbackQuery{callback}, "Action", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track Action %s", answer.Status)
		appMetrika <- true
	})

	<-appMetrika // Send track to Yandex.AppMetrika
}
