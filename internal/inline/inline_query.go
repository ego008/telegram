package inline

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/requests"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func InlineQuery(query *tg.InlineQuery) {
	user, err := database.DB.GetUser(query.From)
	errors.Check(err)

	localizer := i18n.I18N.NewLocalizer(user.Locale, query.From.LanguageCode)

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
	answer.SwitchPrivateMessageText = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "inline_button_dashboard",
	})
	answer.SwitchPrivateMessageParameter = "settings"

	results := getResults(
		query,
		user,
		&requests.Params{
			PageID: offset,
			Tags:   query.Query,
		},
	)

	switch {
	case offset <= 0 &&
		len(results) <= 0:
		answer.SwitchPrivateMessageText = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "inline_no_result",
		})
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

func getResults(query *tg.InlineQuery, user *models.User, p *requests.Params) (results []interface{}) {
	if len(user.WhiteList) > 0 {
		p.Tags += fmt.Sprint(" ", strings.Join(user.WhiteTags(), " "))
	}

	if len(user.BlackList) > 0 {
		p.Tags += fmt.Sprint(" -", strings.Join(user.BlackTags(), " -"))
	}

	filters := user.GetRatingsFilter()
	if filters != "" {
		p.Tags += fmt.Sprint(" ", user.GetRatingsFilter())
	}

	log.Ln("getResults")
	var resNames []string
	for _, res := range user.Resources {
		if resources.Resources[res.Name] == nil {
			continue
		}
		resNames = append(resNames, res.Name)
	}

	if len(resNames) <= 0 {
		return nil
	}

	p.Limit = 50 / len(resNames)

	log.Ln("getResults after resources")
	var wg sync.WaitGroup
	wg.Add(len(resNames))

	log.Ln("getResults preparing res")
	for _, res := range resNames {
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
				result := getResultByPost(query, user, res, &posts[j])
				if result != nil {
					results = append(results, result)
				}
			}
		}(res, p)
	}

	wg.Wait()

	return results
}

func getResultByPost(query *tg.InlineQuery, user *models.User, res string, post *models.Result) interface{} {
	T, err := i18n.SwitchTo(user.Locale, query.From.LanguageCode)
	errors.Check(err)

	replyMarkup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(
				T("button_original"), post.FileURL(resNames).String(),
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
		if !user.ContentTypes.Video {
			return nil
		}

		inputMessageContent := tg.NewInputTextMessageContent(
			localize.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "message_blushboard",
				TemplateData: map[string]string{
					"Type": strings.Title(localize.MustLocalize(&i18n.LocalizeConfig{
						MessageID: "type_video",
					})),
					"Owner": post.Owner,
					"URL": fmt.Sprint(
						resources.Resources[res].UString("scheme", "http"),
						"://",
						resources.Resources[res].UString("host"),
						resources.Resources[res].UString("result"),
						post.ID,
					),
				},
			}),
		)
		inputMessageContent.ParseMode = tg.ModeMarkdown
		inputMessageContent.DisableWebPagePreview = false

		video := tg.NewInlineQueryResultVideo(
			fmt.Sprint(res, post.ID),
			post.FileURL(res).String(),
			tg.MimeHTML,
			post.PreviewURL(res).String(),
			localize.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "inline_title",
				TemplateData: map[string]string{
					"Type": strings.Title(localize.MustLocalize(&i18n.LocalizeConfig{
						MessageID: "type_video",
					})),
					"Owner": post.Owner,
				},
			}),
		)
		video.VideoWidth = post.Width
		video.VideoHeight = post.Height
		video.Description = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "inline_description",
			TemplateData: map[string]string{
				"Rating": post.Rating,
				"Tags":   post.Tags,
			},
		})
		video.InputMessageContent = inputMessageContent
		video.ReplyMarkup = replyMarkup
		return video
	case strings.HasSuffix(post.FileURL(res).Path, "gif"):
		if !user.ContentTypes.Animation {
			return nil
		}

		gif := tg.NewInlineQueryResultGif(
			fmt.Sprint(res, post.ID),
			post.FileURL(res).String(),
			post.FileURL(res).String(),
		)
		gif.GifWidth = post.Width
		gif.GifHeight = post.Height
		gif.Title = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "inline_title",
			TemplateData: map[string]string{
				"Type": strings.Title(localize.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "type_animation",
				})),
				"Owner": post.Owner,
			},
		})
		gif.ReplyMarkup = replyMarkup
		return gif
	default:
		if !user.ContentTypes.Image {
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

		photo.Title = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "inline_title",
			TemplateData: map[string]string{
				"Type": strings.Title(localize.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "type_image",
				})),
				"Owner": post.Owner,
			},
		})
		photo.Description = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "inline_description",
			TemplateData: map[string]string{
				"Rating": post.Rating,
				"Tags":   post.Tags,
			},
		})
		photo.ReplyMarkup = replyMarkup
		return photo
	}
}
