package inline

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/db"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/models"
	"github.com/HentaiDB/HentaiDBot/internal/requests"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func InlineQuery(query *tg.InlineQuery) {
	usr, err := db.GetUserElseAdd(query.From.ID, query.From.LanguageCode)
	errors.Check(err)

	T, err := i18n.SwitchTo(usr.Language, query.From.LanguageCode)
	errors.Check(err)

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
		&requests.Params{
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

	_, err = bot.Bot.AnswerInlineQuery(answer)
	if err != nil {
		log.Ln("[ERROR]", err.Error())
	}
}

func getResults(query *tg.InlineQuery, usr *models.User, p *requests.Params) (results []interface{}) {
	if len(usr.Whitelist) > 0 {
		p.Tags += fmt.Sprint(" ", strings.Join(usr.Whitelist, " "))
	}

	if len(usr.Blacklist) > 0 {
		p.Tags += fmt.Sprint(" -", strings.Join(usr.Blacklist, " -"))
	}

	filters := usr.GetRatingsFilter()
	if filters != "" {
		p.Tags += fmt.Sprint(" ", usr.GetRatingsFilter())
	}

	log.Ln("getResults")
	var res []string
	for key, on := range usr.Resources {
		if on &&
			resources.Resources[key] != nil {
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
		go func(res string, p *requests.Params) {
			defer wg.Done()
			log.Ln("Getted", p.Limit, "results from", res, "...")
			log.D(resources.Resources[res].UMap(""))

			posts, err := requests.Results(res, p)
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

func getResultByPost(query *tg.InlineQuery, usr *models.User, res string, post *models.Result) interface{} {
	T, err := i18n.SwitchTo(usr.Language, query.From.LanguageCode)
	errors.Check(err)

	replyMarkup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(
				T("button_original"), post.FileURL(res).String(),
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
	case strings.HasSuffix(post.FileURL(res).Path, "webm"):
		if !usr.ContentTypes.Video {
			return nil
		}

		inputMessageContent := tg.NewInputTextMessageContent(
			T("message_blushboard", map[string]interface{}{
				"Type":  strings.Title(T("type_video")),
				"Owner": post.Owner,
				"URL": fmt.Sprint(
					resources.Resources[res].UString("scheme", "http"),
					"://",
					resources.Resources[res].UString("host"),
					resources.Resources[res].UString("result"),
					post.ID,
				),
			}),
		)
		inputMessageContent.ParseMode = tg.ModeMarkdown
		inputMessageContent.DisableWebPagePreview = false

		video := tg.NewInlineQueryResultVideo(
			fmt.Sprint(res, post.ID),
			post.FileURL(res).String(),
			tg.MimeHTML,
			post.PreviewURL(res).String(),
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
	case strings.HasSuffix(post.FileURL(res).Path, "gif"):
		if !usr.ContentTypes.Animation {
			return nil
		}

		gif := tg.NewInlineQueryResultGif(
			fmt.Sprint(res, post.ID),
			post.FileURL(res).String(),
			post.FileURL(res).String(),
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
		if !usr.ContentTypes.Image {
			return nil
		}

		photo := tg.NewInlineQueryResultPhoto(
			fmt.Sprint(res, post.ID),
			post.FileURL(res).String(),
			post.PreviewURL(res).String(),
		)
		photo.PhotoWidth = post.Width
		photo.PhotoHeight = post.Height

		if post.Sample {
			photo.PhotoURL = post.SampleURL(res).String()
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
