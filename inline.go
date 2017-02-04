package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/i18n"
)

const BlushBoard = "http://beta.hentaidb.pw"

var (
	results []interface{}
	rating  = map[string]string{
		"s": T("rating_safe"),
		"e": T("rating_explicit"),
		"q": T("rating_questionable"),
		"?": T("rating_unknown"),
	}
)

func getInlineResults(cacheTime int, inline *tg.InlineQuery) {
	// Track action
	b.TrackAsync(inline.From.ID, MetrikaInlineQuery{inline}, "Search", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track Search %s", answer.Status)
		metrika <- true
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
			inlineResult(post[i], T)
		}
	case len(post) == 0: // Found nothing
		emptyKeyboard := tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonURL(T("button_channel"), cfg["link_channel"].(string)),
				tg.NewInlineKeyboardButtonURL(T("button_group"), cfg["link_group"].(string)),
			),
		)
		empty := tg.NewInlineQueryResultArticleMarkdown(inline.ID, T("inline_no_result_title"), "`¯\\_(ツ)_/¯`")
		empty.Description = T("inline_no_result_description")
		empty.ReplyMarkup = &emptyKeyboard
		results = append(results, empty)
	}

	// Configure inline-mode
	var inlineConfig tg.InlineConfig
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

	<-metrika // Send track to Yandex.AppMetrika
}

func inlineResult(post Post, locale i18n.TranslateFunc) {
	// Universal(?) preview url
	preview := fmt.Sprint(cfg["resource_url"].(string), cfg["resource_thumbs_dir"].(string), post.Directory, cfg["resource_thumbs_part"].(string), post.Hash, ".jpg")
	post.Rating = rating[post.Rating]
	resultKeyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(T("button_original"), "https:"+post.FileURL),
		),
	)

	switch {
	case strings.Contains(post.FileURL, ".webm"): // Not support yet. Show tip about "hidden" result
		BBURL := BlushBoard + "/hash/" + post.Hash

		video := tg.NewInlineQueryResultVideo(strconv.Itoa(post.ID), BlushBoard+"/embed/"+post.Hash)
		video.MimeType = "text/html"
		video.ThumbURL = preview
		video.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_video")),
			"Owner": post.Owner,
		})
		video.Width = post.Width
		video.Height = post.Height
		video.Description = T("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		video.InputMessageContent = tg.InputTextMessageContent{
			Text: T("message_blushboard", map[string]interface{}{
				"Type":  strings.Title(T("type_video")),
				"Owner": post.Owner,
				"URL":   BBURL,
			}),
			ParseMode:             parseMarkdown,
			DisableWebPagePreview: false,
		}
		results = append(results, video)
	case strings.Contains(post.FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
		video := tg.NewInlineQueryResultVideo(strconv.Itoa(post.ID), "https:"+post.FileURL)
		video.MimeType = "video/mp4"
		video.ThumbURL = preview
		video.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_video")),
			"Owner": post.Owner,
		})
		video.Width = post.Width
		video.Height = post.Height
		video.Description = T("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		video.ReplyMarkup = &resultKeyboard
		results = append(results, video)
	case strings.Contains(post.FileURL, ".gif"):
		gif := tg.NewInlineQueryResultGIF(strconv.Itoa(post.ID), "https:"+post.FileURL)
		gif.Width = post.Width
		gif.Height = post.Height
		gif.ThumbURL = post.FileURL
		gif.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_animation")),
			"Owner": post.Owner,
		})
		gif.ReplyMarkup = &resultKeyboard
		results = append(results, gif)
	default:
		image := tg.NewInlineQueryResultPhoto(strconv.Itoa(post.ID), "https:"+post.FileURL)
		image.ThumbURL = preview
		image.Width = post.Width
		image.Height = post.Height
		image.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_image")),
			"Owner": post.Owner,
		})
		image.Description = T("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		image.ReplyMarkup = &resultKeyboard
		results = append(results, image)
	}
}

func trackInlineResult(result *tg.ChosenInlineResult) {
	b.TrackAsync(result.From.ID, MetrikaChosenInlineResult{result}, "Find", func(answer botan.Answer, err []error) {
		log.Printf("[Botan] Track Find %s", answer.Status)
		metrika <- true
	})

	<-metrika // Send track to Yandex.AppMetrika
}
