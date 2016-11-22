package main

import (
	b "github.com/botanio/sdk/go"
	t "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"strconv"
	"strings"
)

var results []interface{}

const BlushBoard = "http://beta.hentaidb.pw"

func getInlineResults(cacheTime int, inline *t.InlineQuery) {
	locale := checkLanguage(inline.From)

	// Track action
	metrika.TrackAsync(inline.From.ID, MetrikaInlineQuery{inline}, "Search", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track Search %s", answer.Status)
		appMetrika <- true
	})

	nsfw := checkNSFW(inline.From)
	if !nsfw {
		inline.Query += " rating:safe"
	}

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
			inlineResult(post[i], locale)
		}
	case len(post) == 0: // Found nothing
		emptyKeyboard := t.NewInlineKeyboardMarkup(
			t.NewInlineKeyboardRow(
				t.NewInlineKeyboardButtonURL(locale("button_channel"), config.Links.Channel),
				t.NewInlineKeyboardButtonURL(locale("button_group"), config.Links.Group),
			),
		)
		empty := t.NewInlineQueryResultArticleMarkdown(inline.ID, locale("inline_no_result_title"), "`¯\\_(ツ)_/¯`")
		empty.Description = locale("inline_no_result_description")
		empty.ReplyMarkup = &emptyKeyboard
		results = append(results, empty)
	}

	// Configure inline-mode
	inlineConfig := t.InlineConfig{}
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

func inlineResult(post Post, locale i18n.TranslateFunc) {
	// Universal(?) preview url
	preview := config.Resource[20].Settings.URL + config.Resource[20].Settings.ThumbsDir + post.Directory + config.Resource[20].Settings.ThumbsPart + post.Hash + ".jpg"
	post.Rating = setResultRating(post.Rating, locale)
	resultKeyboard := t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonURL(locale("button_original"), post.FileURL),
		),
	)

	switch {
	case strings.Contains(post.FileURL, ".webm"): // Not support yet. Show tip about "hidden" result
		BBURL := BlushBoard + "/hash/" + post.Hash

		video := t.NewInlineQueryResultVideo(strconv.Itoa(post.ID), BlushBoard+"/embed/"+post.Hash)
		video.MimeType = "text/html"
		video.ThumbURL = preview
		video.Title = locale("inline_title", map[string]interface{}{
			"Type":  strings.Title(locale("type_video")),
			"Owner": post.Owner,
		})
		video.Width = post.Width
		video.Height = post.Height
		video.Description = locale("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		video.InputMessageContent = t.InputTextMessageContent{
			Text: locale("message_blushboard", map[string]interface{}{
				"Type":  strings.Title(locale("type_video")),
				"Owner": post.Owner,
				"URL":   BBURL,
			}),
			ParseMode:             parseMarkdown,
			DisableWebPagePreview: false,
		}
		results = append(results, video)
	case strings.Contains(post.FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
		video := t.NewInlineQueryResultVideo(strconv.Itoa(post.ID), post.FileURL)
		video.MimeType = "video/mp4"
		video.ThumbURL = preview
		video.Title = locale("inline_title", map[string]interface{}{
			"Type":  strings.Title(locale("type_video")),
			"Owner": post.Owner,
		})
		video.Width = post.Width
		video.Height = post.Height
		video.Description = locale("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		video.ReplyMarkup = &resultKeyboard
		results = append(results, video)
	case strings.Contains(post.FileURL, ".gif"):
		gif := t.NewInlineQueryResultGIF(strconv.Itoa(post.ID), post.FileURL)
		gif.Width = post.Width
		gif.Height = post.Height
		gif.ThumbURL = post.FileURL
		gif.Title = locale("inline_title", map[string]interface{}{
			"Type":  strings.Title(locale("type_animation")),
			"Owner": post.Owner,
		})
		gif.ReplyMarkup = &resultKeyboard
		results = append(results, gif)
	default:
		image := t.NewInlineQueryResultPhoto(strconv.Itoa(post.ID), post.FileURL)
		image.ThumbURL = preview
		image.Width = post.Width
		image.Height = post.Height
		image.Title = locale("inline_title", map[string]interface{}{
			"Type":  strings.Title(locale("type_image")),
			"Owner": post.Owner,
		})
		image.Description = locale("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		image.ReplyMarkup = &resultKeyboard
		results = append(results, image)
	}
}

func setResultRating(rating string, locale i18n.TranslateFunc) string {
	switch rating {
	case "s":
		return locale("rating_safe")
	case "e":
		return locale("rating_explicit")
	case "q":
		return locale("rating_questionable")
	default:
		return locale("rating_unknown")
	}
}

func trackInlineResult(result *t.ChosenInlineResult) {
	metrika.TrackAsync(result.From.ID, MetrikaChosenInlineResult{result}, "Find", func(answer b.Answer, err []error) {
		log.Printf("[Botan] Track Find %s", answer.Status)
		appMetrika <- true
	})

	<-appMetrika // Send track to Yandex.AppMetrika
}
