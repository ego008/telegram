package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	i18n "github.com/nicksnyder/go-i18n/i18n"
)

var exceptions = map[string]bool{
	"height":  true,
	"id":      true,
	"parent":  true,
	"rating":  true,
	"score":   true,
	"source":  true,
	"updated": true,
	"user":    true,
	"width":   true,
	"md5":     true,
	"sort":    true,
}

func inline(inline *tg.InlineQuery) {
	trackInline(inline)

	usr, err := getUser(inline.From.ID)
	if err != nil {
		log.Println("Create user:", err.Error())
	}

	T, err := i18n.Tfunc(usr.Language)
	if err != nil {
		log.Println(err.Error())
	}

	inline.Query = strings.ToLower(inline.Query)

	var page int
	hasRating := false
	if r.MatchString(inline.Query) {
		operators := r.FindAllStringSubmatch(inline.Query, -1)
		for _, operator := range operators {
			op := strings.Split(operator[1], ":")
			if exceptions[op[0]] {
				if op[0] == "rating" {
					hasRating = true
				}
				continue
			}

			switch op[0] {
			case "page":
				page, _ = strconv.Atoi(op[1])
			case "lang":
				T, _ = i18n.Tfunc(op[1])
			}

			inline.Query = strings.TrimSuffix(strings.Replace(inline.Query, operator[0], "", -1), " ")
		}
	}

	// Check result pages
	var posts []gPost
	req := &request{Limit: 50}
	req.Tags = inline.Query

	if len(usr.Whitelist) > 0 {
		req.Tags += fmt.Sprint(" ", strings.Join(usr.Whitelist, " "))
	}

	if len(usr.Blacklist) > 0 {
		req.Tags += fmt.Sprint(" -", strings.Join(usr.Blacklist, " -"))
	}

	rt := usr.Ratings
	if !hasRating {
		switch {
		case rt.Safe && !rt.Questionale && !rt.Explicit:
			req.Tags += " rating:safe"
		case !rt.Safe && rt.Questionale && !rt.Explicit:
			req.Tags += " rating:questionable"
		case !rt.Safe && !rt.Questionale && rt.Explicit:
			req.Tags += " rating:explicit"
		case !rt.Safe && rt.Questionale && rt.Explicit:
			req.Tags += " -rating:safe"
		case rt.Safe && !rt.Questionale && rt.Explicit:
			req.Tags += " -rating:questionable"
		case rt.Safe && rt.Questionale && !rt.Explicit:
			req.Tags += " -rating:explicit"
		}
	}

	switch {
	case len(inline.Offset) <= 0:
		req.PageID = page
		posts, _ = getPosts(req)
	case len(inline.Offset) > 0:
		page, err = strconv.Atoi(inline.Offset)
		if err != nil {
			log.Println(err.Error())
		}
		req.PageID = page
		posts, _ = getPosts(req)
	}

	results := collectResults(usr, inline, T, posts)

	// Configure inline-mode
	inlineCfg := tg.InlineConfig{
		CacheTime:     *flagCache,
		InlineQueryID: inline.ID,
		IsPersonal:    true,
		Results:       results,
	}

	switch {
	case len(posts) <= 0:
		inlineCfg.SwitchPMParameter = cheatsheet
		inlineCfg.SwitchPMText = T("inline_no_result")
	case len(posts)%50 == 0:
		page++
		inlineCfg.NextOffset = strconv.Itoa(page)
		inlineCfg.SwitchPMParameter = settings
		inlineCfg.SwitchPMText = T("inline_button_dashboard")
	default:
		inlineCfg.SwitchPMParameter = settings
		inlineCfg.SwitchPMText = T("inline_button_dashboard")
	}

	if _, err := bot.AnswerInlineQuery(inlineCfg); err != nil {
		log.Println("AnswerInlineQuery:", err.Error())
	}

	<-appMetrika // Send track to Yandex.AppMetrika
}

func collectResults(usr *User, inline *tg.InlineQuery, T i18n.TranslateFunc, posts []gPost) []interface{} {
	var results []interface{}
	if len(posts) > 0 {
		for _, post := range posts {
			post.FileURL = fmt.Sprint("https:", post.FileURL)

			if len(post.Tags) >= 30 {
				post.Tags = fmt.Sprint(post.Tags[:30], "...")
			}

			preview := fmt.Sprint(
				glbr["url"].(string),
				glbr["thumbs_dir"].(string),
				post.Directory,
				glbr["thumbs_part"].(string),
				post.Hash,
				".jpg",
			)

			markup := tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButtonURL(T("button_original"), post.FileURL),
				),
			)

			if !strings.HasSuffix(post.Image, ".webm") {
				markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], tg.NewInlineKeyboardButtonData(
					T("button_info"),
					fmt.Sprintf("info %s", fmt.Sprint("glbr", post.ID)),
				))
			}

			if inline.Query != "" {
				markup.InlineKeyboard[0] = append(markup.InlineKeyboard[0], tg.InlineKeyboardButton{
					Text: T("button_more"),
					SwitchInlineQueryCurrentChat: &inline.Query,
				})
			}

			switch post.Rating {
			case "s":
				post.Rating = T("rating_safe")
			case "q":
				post.Rating = T("rating_questionable")
			case "e":
				post.Rating = T("rating_explicit")
			default:
				post.Rating = T("rating_unknown")
			}

			switch {
			case strings.Contains(post.FileURL, ".webm"): // Not support yet. Show tip about "hidden" result
				video := tg.NewInlineQueryResultVideo(fmt.Sprint("glbr", post.ID), fmt.Sprint(bb, "/embed/", post.Hash))
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
						"URL":   fmt.Sprint(bb, "/hash/", post.Hash),
					}),
					ParseMode:             tg.ModeMarkdown,
					DisableWebPagePreview: false,
				}
				results = append(results, video)
			case strings.Contains(post.FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
				video := tg.NewInlineQueryResultVideo(fmt.Sprint("glbr", post.ID), post.FileURL)
				video.MimeType = "video/mp4"
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
				video.ReplyMarkup = &markup
				results = append(results, video)
			case strings.Contains(post.FileURL, ".gif"):
				gif := tg.NewInlineQueryResultGIF(fmt.Sprint("glbr", post.ID), post.FileURL)
				gif.Width = post.Width
				gif.Height = post.Height
				gif.ThumbURL = post.FileURL
				gif.Title = T("inline_title", map[string]interface{}{
					"Type":  strings.Title(T("type_animation")),
					"Owner": post.Owner,
				})
				gif.ReplyMarkup = &markup
				results = append(results, gif)
			default:
				image := tg.NewInlineQueryResultPhoto(fmt.Sprint("glbr", post.ID), post.FileURL)
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
				image.ReplyMarkup = &markup
				results = append(results, image)
			}
		}
	}
	return results
}

func chosenResult(result *tg.ChosenInlineResult) {
	trackChosenResult(result)
	<-appMetrika // Send track to Yandex.AppMetrika
}
