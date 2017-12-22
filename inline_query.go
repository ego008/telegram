package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/go-telegram"
)

var urlIV = &url.URL{
	Scheme: "https",
	Host:   "t.me",
	Path:   "/iv",
}

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
	case len(results) <= 0:
		answer.SwitchPrivateMessageText = T("inline_no_result")
		answer.SwitchPrivateMessageParameter = "settings"
	case offset > 0 && len(results) <= 0:
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

	p.Limit = 50 / len(res)

	log.Ln("getResults after resources")
	var wg sync.WaitGroup
	wg.Add(len(res))

	log.Ln("getResults preparing res")
	for i := range res {
		go func(res string, p *params) {
			defer wg.Done()
			log.Ln("Getted", p.Limit, "results from", res, "...")
			posts, err := request(res, p)
			if err != nil {
				log.Ln("[ERROR]", err.Error())
				return
			}

			log.Ln("Getted", len(posts), "results from", res)
			for j := range posts {
				results = append(results, getResultByPost(query, usr, res, &posts[j]))
			}
		}(res[i], p)
	}

	wg.Wait()

	return results
}

func getResultByPost(query *tg.InlineQuery, usr *user, res string, post *post) interface{} {
	log.Ln("getResultByPost")
	T, err := langSwitch(usr.Language, query.From.LanguageCode)
	errCheck(err)

	resource := resources[res]
	id := fmt.Sprint(res, post.ID)

	urlScheme := checkInterface(resource["scheme"])
	urlHost := checkInterface(resource["host"])

	urlImagesDir := checkInterface(resource["images_dir"])
	urlImagesPart := checkInterface(resource["images_part"])

	urlThumbsDir := checkInterface(resource["thumbs_dir"])
	urlThumbsPart := checkInterface(resource["thumbs_part"])

	file := strings.Split(post.Image, ".")
	fileHash := file[0]
	urlFileFormat := file[1]

	urlFileHash := fileHash
	urlThumbsHash := post.Hash

	if resource["hash"] != nil {
		log.Ln("Hash by:", resource["hash"].(string))
		switch resource["hash"].(string) {
		case "file":
			urlThumbsHash = urlFileHash
		case "thumb":
			urlFileHash = urlThumbsHash
		case "invert":
			urlFileHash = post.Hash
			urlThumbsHash = fileHash
		}
	}

	fileURL := &url.URL{
		Scheme: urlScheme,
		Host:   urlHost,
		Path: fmt.Sprint(
			urlImagesDir,
			post.Directory,
			urlImagesPart, urlFileHash, ".", urlFileFormat,
		),
	}

	urlThumbsFormat := "jpg"
	if resource["thumbs_format"] != nil {
		urlThumbsFormat = resource["thumbs_format"].(string)
		if urlThumbsFormat == "auto" {
			urlThumbsFormat = urlFileFormat
		}
	}

	thumbURL := &url.URL{
		Scheme: urlScheme,
		Host:   urlHost,
		Path: fmt.Sprint(
			urlThumbsDir, post.Directory,
			urlThumbsPart, urlThumbsHash, ".", urlThumbsFormat,
		),
	}

	log.Ln("PreviewURL:", thumbURL.String())
	log.Ln("FileURL:", fileURL.String())

	replyMarkup := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(
				T("button_original"), fileURL.String(),
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

	switch urlFileFormat {
	case "webm":
		inputMessageContent := tg.NewInputTextMessageContent(
			T("message_blushboard", map[string]interface{}{
				"Type":  strings.Title(T("type_video")),
				"Owner": post.Owner,
				"URL": fmt.Sprint(
					urlScheme, "://",
					urlHost,
					checkInterface(resource["result"]), post.ID,
				),
			}),
		)
		inputMessageContent.ParseMode = tg.ModeMarkdown
		inputMessageContent.DisableWebPagePreview = false

		video := tg.NewInlineQueryResultVideo(
			id,
			fileURL.String(),
			tg.MimeHTML,
			thumbURL.String(),
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
	case "gif":
		gif := tg.NewInlineQueryResultGif(id, fileURL.String(), fileURL.String())
		gif.GifWidth = post.Width
		gif.GifHeight = post.Height
		gif.Title = T("inline_title", map[string]interface{}{
			"Type":  strings.Title(T("type_animation")),
			"Owner": post.Owner,
		})
		gif.ReplyMarkup = replyMarkup
		return gif
	default:
		photo := tg.NewInlineQueryResultPhoto(id, fileURL.String(), thumbURL.String())
		photo.PhotoWidth = post.Width
		photo.PhotoHeight = post.Height
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
