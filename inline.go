package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

func getInlineResults(inline *tgbotapi.InlineQuery, cacheTime int) {
	// Check result pages
	var posts []Post
	var resultPage int
	if len(update.InlineQuery.Offset) > 0 {
		resultPage, _ = strconv.Atoi(update.InlineQuery.Offset)
		posts = getPosts(Request{PageID: resultPage, Tags: update.InlineQuery.Query})
	} else {
		posts = getPosts(Request{Tags: update.InlineQuery.Query})
	}

	// Analysis of results
	var result []interface{}
	switch {
	case len(posts) > 0:
		for i := 0; i < len(posts); i++ {
			// Universal(?) preview url
			preview := config.Resource[resNum].Settings.URL + config.Resource[resNum].Settings.ThumbsDir + posts[i].Directory + config.Resource[resNum].Settings.ThumbsPart + posts[i].Hash + ".jpg"

			// Rating
			var rating string
			switch posts[i].Rating {
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
					tgbotapi.NewInlineKeyboardButtonURL("Original image", posts[i].FileURL),
				),
			)

			switch {
			case strings.Contains(posts[i].FileURL, ".webm"): // It is necessary to get around error 403 when requesting video :|
				// query := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(i), posts[i].FileURL) // Does not work
				// query.MimeType = "text/html" // Link on widget-page?
				// query.MimeType = "video/mp4" // Does not work for .webm
				// query.ThumbURL = preview
				// query.Width = posts[i].Width
				// query.Height = posts[i].Height
				// query.Title = "Video by " + strings.Title(posts[i].Owner)
				// query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
				// query.ReplyMarkup = &button
				// result = append(result, query)
				continue
			case strings.Contains(posts[i].FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
				query := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(i), posts[i].FileURL)
				query.MimeType = "video/mp4"
				query.ThumbURL = preview
				query.Width = posts[i].Width
				query.Height = posts[i].Height
				query.Title = "Video by " + strings.Title(posts[i].Owner)
				query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
				query.ReplyMarkup = &button
				result = append(result, query)
			case strings.Contains(posts[i].FileURL, ".gif"):
				query := tgbotapi.NewInlineQueryResultGIF(strconv.Itoa(i), posts[i].FileURL)
				query.ThumbURL = posts[i].FileURL
				query.Width = posts[i].Width
				query.Height = posts[i].Height
				query.Title = "Animation by " + strings.Title(posts[i].Owner)
				query.ReplyMarkup = &button
				result = append(result, query)
			default:
				query := tgbotapi.NewInlineQueryResultPhoto(strconv.Itoa(i), posts[i].FileURL)
				query.ThumbURL = preview
				query.Width = posts[i].Width
				query.Height = posts[i].Height
				query.Title = "Image by " + strings.Title(posts[i].Owner)
				query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
				query.ReplyMarkup = &button
				result = append(result, query)
			}
		}
	case len(posts) == 0: // Found nothing
		query := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, noInlineResultTitle, noInlineResultMessage)
		query.Description = noInlineResultDescription
		result = append(result, query)
	}

	// Configure inline-mode
	inlineConfig := tgbotapi.InlineConfig{}
	inlineConfig.InlineQueryID = update.InlineQuery.ID
	inlineConfig.IsPersonal = true
	inlineConfig.CacheTime = cacheTime
	inlineConfig.Results = result
	if len(posts) == 50 {
		inlineConfig.NextOffset = strconv.Itoa(resultPage + 1) // If available next page of results
	}

	if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
		log.Printf("[Bot] Answer inline-query error: %+v", err)
	}
}
