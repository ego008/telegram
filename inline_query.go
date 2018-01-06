package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func inlineQuery(query *tg.InlineQuery) {
	usr, err := dbGetUserElseAdd(query.From.ID, query.From.LanguageCode)
	errCheck(err)

	T, err := langSwitch(usr.Language, query.From.LanguageCode)
	errCheck(err)

	if query.Offset == "" {
		query.Offset = "-1"
	}

	offset, err := strconv.Atoi(query.Offset)
	if err != nil {
		log.Ln("[ERROR]", err.Error())
	}
	offset++

	answer := tg.NewAnswerInlineQuery(query.ID)
	answer.CacheTime = 1
	answer.IsPersonal = true
	answer.SwitchPrivateMessageText = T("inline_button_dashboard")
	answer.SwitchPrivateMessageParameter = "settings"

	results := getResults(
		query,
		usr,
		&params{
			PageID: offset,
			Tags:   query.Query,
		},
	)

	switch {
	case offset <= 0 &&
		len(results) <= 0:
		answer.SwitchPrivateMessageText = T("inline_no_result")
		answer.SwitchPrivateMessageParameter = "settings"
	case offset > 0 &&
		len(results) <= 0:
	default:
		answer.Results = results
		answer.NextOffset = strconv.Itoa(offset)
	}

	_, err = bot.AnswerInlineQuery(answer)
	if err != nil {
		log.Ln("[ERROR]", err.Error())
	}
}

func getResults(query *tg.InlineQuery, usr *user, p *params) (results []interface{}) {
	if len(usr.Whitelist) > 0 {
		p.Tags += fmt.Sprint(" ", strings.Join(usr.Whitelist, " "))
	}

	if len(usr.Blacklist) > 0 {
		p.Tags += fmt.Sprint(" -", strings.Join(usr.Blacklist, " -"))
	}

	filters := usr.getRatingsFilter()
	if filters != "" {
		p.Tags += fmt.Sprint(" ", usr.getRatingsFilter())
	}

	log.Ln("getResults")
	var res []string
	for key, on := range usr.Resources {
		if on &&
			resources[key] != nil {
			res = append(res, key)
		}
	}

	if len(res) <= 0 {
		return nil
	}

	p.Limit = 50 / len(res)

	log.Ln("getResults after resources")
	var wg sync.WaitGroup
	wg.Add(len(res))

	log.Ln("getResults preparing res")
	for i := range res {
		go func(res string, p *params) {
			defer wg.Done()
			log.Ln("Getted", p.Limit, "results from", res, "...")
			log.D(resources[res].UMap(""))

			posts, err := request(res, p)
			if err != nil {
				log.Ln("[ERROR]", err.Error())
				return
			}

			log.Ln("Getted", len(posts), "results from", res)
			for j := range posts {
				result := getResultByPost(query, usr, res, &posts[j])
				if result != nil {
					results = append(results, result)
				}
			}
		}(res[i], p)
	}

	wg.Wait()

	return results
}

func getResultByPost(query *tg.InlineQuery, usr *user, res string, post *post) interface{} {
	T, err := langSwitch(usr.Language, query.From.LanguageCode)
	errCheck(err)

	replyMarkup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(
				T("button_original"), post.fileURL(res).String(),
			),
		),
	)

	if query.Query != "" {
		replyMarkup.InlineKeyboard[0] = append(
			replyMarkup.InlineKeyboard[0],
			tg.NewInlineKeyboardButtonSwitchSelf(
				T("button_more"), query.Query,
			),
		)
	}

	switch {
	case strings.HasSuffix(post.fileURL(res).Path, "webm"):
		if !usr.Types.Video {
			return nil
		}

		inputMessageContent := tg.NewInputTextMessageContent(
			T("message_blushboard", map[string]interface{}{
				"Type":  strings.Title(T("type_video")),
				"Owner": post.Owner,
				"URL": fmt.Sprint(
					resources[res].UString("scheme", "http"),
					"://",
					resources[res].UString("host"),
					resources[res].UString("result"),
					post.ID,
				),
			}),
		)
		inputMessageContent.ParseMode = tg.ModeMarkdown
		inputMessageContent.DisableWebPagePreview = false

		video := tg.NewInlineQueryResultVideo(
			fmt.Sprint(res, post.ID),
			post.fileURL(res).String(),
			tg.MimeHTML,
			post.previewURL(res).String(),
			T("inline_title", map[string]interface{}{
				"Type":  strings.Title(T("type_video")),
				"Owner": post.Owner,
			}),
		)
		video.VideoWidth = post.Width
		video.VideoHeight = post.Height
		video.Description = T("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		video.InputMessageContent = inputMessageContent
		video.ReplyMarkup = replyMarkup
		return video
	case strings.HasSuffix(post.fileURL(res).Path, "gif"):
		if !usr.Types.Animation {
			return nil
		}

		gif := tg.NewInlineQueryResultGif(
			fmt.Sprint(res, post.ID),
			post.fileURL(res).String(),
			post.fileURL(res).String(),
		)
		gif.GifWidth = post.Width
		gif.GifHeight = post.Height
		gif.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_animation")),
			"Owner": post.Owner,
		})
		gif.ReplyMarkup = replyMarkup
		return gif
	default:
		if !usr.Types.Image {
			return nil
		}

		photo := tg.NewInlineQueryResultPhoto(
			fmt.Sprint(res, post.ID),
			post.fileURL(res).String(),
			post.previewURL(res).String(),
		)
		photo.PhotoWidth = post.Width
		photo.PhotoHeight = post.Height

		if post.Sample {
			photo.PhotoURL = post.sampleURL(res).String()
			photo.PhotoWidth = post.SampleWidth
			photo.PhotoHeight = post.SampleHeight
		}

		photo.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_image")),
			"Owner": post.Owner,
		})
		photo.Description = T("inline_description", map[string]interface{}{
			"Rating": post.Rating,
			"Tags":   post.Tags,
		})
		photo.ReplyMarkup = replyMarkup
		return photo
	}
}
