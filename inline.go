package main

import (
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
				rating = "Safe"
			case "e":
				rating = "Explicit"
			case "q":
				rating = "Questionable"
			default:
				rating = "Unknown"
			}

			button := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Original image", post[i].FileURL),
				),
			)

			switch {
			case strings.Contains(post[i].FileURL, ".webm"): // It is necessary to get around error 403 when requesting video :|
				// video := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(i), post[i].FileURL) // Does not work
				// video.MimeType = "text/html" // Link on widget-page?
				// video.MimeType = "video/mp4" // Does not work for .webm
				// video.ThumbURL = preview
				// video.Width = post[i].Width
				// video.Height = post[i].Height
				// video.Title = "Video by " + strings.Title(post[i].Owner)
				// video.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(post[i].Score) + "\nTags: " + post[i].Tags
				// video.ReplyMarkup = &button
				// result = append(result, video)
				continue
			case strings.Contains(post[i].FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
				video := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(i), post[i].FileURL)
				video.MimeType = "video/mp4"
				video.ThumbURL = preview
				video.Width = post[i].Width
				video.Height = post[i].Height
				video.Title = "Video by " + strings.Title(post[i].Owner)
				video.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(post[i].Score) + "\nTags: " + post[i].Tags
				video.ReplyMarkup = &button
				result = append(result, video)
			case strings.Contains(post[i].FileURL, ".gif"):
				gif := tgbotapi.NewInlineQueryResultGIF(strconv.Itoa(i), post[i].FileURL)
				gif.ThumbURL = post[i].FileURL
				gif.Width = post[i].Width
				gif.Height = post[i].Height
				gif.Title = "Animation by " + strings.Title(post[i].Owner)
				gif.ReplyMarkup = &button
				result = append(result, gif)
			default:
				image := tgbotapi.NewInlineQueryResultPhoto(strconv.Itoa(i), post[i].FileURL)
				image.ThumbURL = preview
				image.Width = post[i].Width
				image.Height = post[i].Height
				image.Title = "Image by " + strings.Title(post[i].Owner)
				image.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(post[i].Score) + "\nTags: " + post[i].Tags
				image.ReplyMarkup = &button
				result = append(result, image)
			}
		}
	case len(post) == 0: // Found nothing
		empty := tgbotapi.NewInlineQueryResultArticle(inline.ID, noInlineResultTitle, noInlineResultMessage)
		empty.Description = noInlineResultDescription
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
