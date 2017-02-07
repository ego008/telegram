package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/botanio/sdk/go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/kirillDanshin/dlog"
	"github.com/nicksnyder/go-i18n/i18n"
)

const BlushBoard = "http://beta.hentaidb.pw"

var results []interface{}

func GetInlineResults(inline *tg.InlineQuery) {
	// Track action
	b.TrackAsync(inline.From.ID, struct{ *tg.InlineQuery }{inline}, "Search", func(answer botan.Answer, err []error) {
		log.Ln("Track Search", answer.Status)
		metrika <- true
	})

	usr, err := GetUserDB(inline.From.ID)
	if err != nil {
		log.Ln(err.Error())
		return
	}
	T, _ := i18n.Tfunc(usr.Language)

	log.Ln(usr)

	more := inline.Query
	if !usr.NSFW {
		if len(inline.Query) > 0 {
			inline.Query += " rating:safe"
			more = strings.TrimSuffix(more, " rating:safe")
		} else {
			inline.Query = "rating:safe"
			more = strings.TrimSuffix(more, "rating:safe")
		}
	}

	// Check result pages
	var posts []Post
	var page int
	switch {
	case len(inline.Offset) > 0:
		page, _ = strconv.Atoi(inline.Offset)
		posts = getPosts(Request{
			Limit:  50,
			PageID: page,
			Tags:   inline.Query,
		})
	case len(inline.Offset) <= 0:
		posts = getPosts(Request{
			Limit: 50,
			Tags:  inline.Query,
		})
	}

	results = nil
	switch {
	case len(posts) > 0:
		for _, post := range posts {
			inlineResult(more, post, T)
		}
	case len(posts) == 0: // Found nothing
		empty := tg.NewInlineQueryResultArticleMarkdown(inline.ID, T("inline_no_result_title"), "`¯\\_(ツ)_/¯`")
		empty.Description = T("inline_no_result_description")
		markup := tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonURL(T("button_channel"), cfg["link_channel"].(string)),
				tg.NewInlineKeyboardButtonURL(T("button_group"), cfg["link_group"].(string)),
			),
		)
		empty.ReplyMarkup = &markup
		results = append(results, empty)
	}

	// Configure inline-mode
	var inlineConfig tg.InlineConfig
	inlineConfig.SwitchPMText = T("inline_button_dashboard")
	inlineConfig.SwitchPMParameter = "settings"
	inlineConfig.InlineQueryID = inline.ID
	inlineConfig.IsPersonal = true
	inlineConfig.CacheTime = *cacheTime
	inlineConfig.Results = results
	if len(posts) == 50 {
		page++
		inlineConfig.NextOffset = strconv.Itoa(page)
	}

	if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
		log.Ln("Answer inline-query error:", err.Error())
	}

	<-metrika // Send track to Yandex.AppMetrika
}

func inlineResult(query string, post Post, T i18n.TranslateFunc) {
	preview := fmt.Sprint(cfg["resource_url"].(string), cfg["resource_thumbs_dir"].(string), post.Directory, cfg["resource_thumbs_part"].(string), post.Hash, ".jpg")

	markup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(T("button_original"), fmt.Sprint("https:", post.FileURL)),
		),
	)
	if query != "" {
		markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], tg.InlineKeyboardButton{
			Text: T("button_more"),
			SwitchInlineQueryCurrentChat: &query,
		})
	}

	rating := map[string]string{
		"s": T("rating_safe"),
		"e": T("rating_explicit"),
		"q": T("rating_questionable"),
		"?": T("rating_unknown"),
	}
	post.Rating = rating[post.Rating]

	switch {
	case strings.Contains(fmt.Sprint("https:", post.FileURL), ".webm"): // Not support yet. Show tip about "hidden" result
		video := tg.NewInlineQueryResultVideo(strconv.Itoa(post.ID), fmt.Sprintf("%s/embed/%s", BlushBoard, post.Hash))
		video.MimeType = "text/html"
		video.ThumbURL = preview
		video.Width = post.Width
		video.Height = post.Height
		video.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_video")),
			"Owner": post.Owner,
		})
		video.Description = T("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		video.InputMessageContent = tg.InputTextMessageContent{
			Text: T("message_blushboard", map[string]interface{}{
				"Type":  strings.Title(T("type_video")),
				"Owner": post.Owner,
				"URL":   fmt.Sprintf("%s/hash/%s", BlushBoard, post.Hash),
			}),
			ParseMode:             parseMarkdown,
			DisableWebPagePreview: false,
		}
		results = append(results, video)
	case strings.Contains(fmt.Sprint("https:", post.FileURL), ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
		video := tg.NewInlineQueryResultVideo(strconv.Itoa(post.ID), fmt.Sprint("https:", post.FileURL))
		video.MimeType = "video/mp4"
		video.ThumbURL = preview
		video.Width = post.Width
		video.Height = post.Height
		/*
			video.Title = T("inline_title", map[string]interface{}{
				"Type":  strings.Title(T("type_video")),
				"Owner": post.Owner,
			})
			video.Description = T("inline_description", map[string]interface{}{
				"Rating": post.Rating,
				"Tags":   post.Tags,
			})
		*/
		video.ReplyMarkup = &markup
		results = append(results, video)
	case strings.Contains(fmt.Sprint("https:", post.FileURL), ".gif"):
		gif := tg.NewInlineQueryResultGIF(strconv.Itoa(post.ID), fmt.Sprint("https:", post.FileURL))
		gif.Width = post.Width
		gif.Height = post.Height
		gif.ThumbURL = fmt.Sprint("https:", post.FileURL)
		/*
			gif.Title = T("inline_title", map[string]interface{}{
				"Type":  strings.Title(T("type_animation")),
				"Owner": post.Owner,
			})
		*/
		gif.ReplyMarkup = &markup
		results = append(results, gif)
	default:
		image := tg.NewInlineQueryResultPhoto(strconv.Itoa(post.ID), fmt.Sprint("https:", post.FileURL))
		image.ThumbURL = preview
		image.Width = post.Width
		image.Height = post.Height
		/*
			image.Title = T("inline_title", map[string]interface{}{
				"Type":  strings.Title(T("type_image")),
				"Owner": post.Owner,
			})
			image.Description = T("inline_description", map[string]interface{}{
				"Rating": post.Rating,
				"Tags":   post.Tags,
			})
		*/
		image.ReplyMarkup = &markup
		results = append(results, image)
	}
}

func TrackChosenInlineResult(result *tg.ChosenInlineResult) {
	b.TrackAsync(result.From.ID, struct{ *tg.ChosenInlineResult }{result}, "Find", func(answer botan.Answer, err []error) {
		log.Ln("Track Find", answer.Status)
		metrika <- true
	})

	go AddHitsDB(result.From.ID)

	<-metrika // Send track to Yandex.AppMetrika
}
