package commands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	"github.com/HentaiDB/HentaiDBot/internal/requests"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func commandRandom(msg *tg.Message) {
	user, err := database.DataBase.GetUser(msg.From)
	errors.Check(err)

	_, err := bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionUploadPhoto)
	errors.Check(err)

	// T, err := langSwitch(user.Locale, msg.From.LanguageCode)
	// errors.Check(err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := resources.Tags[r.Intn(len(resources.Tags)-1)]

	posts, err := requests.Results(res, &requests.Params{
		Tags: msg.CommandArgument(),
	})
	if err != nil {
		commandRandom(msg)
		return
	}

	if len(posts) <= 0 {
		text := fmt.Sprint("No results by ", msg.CommandArgument(), " tags.")
		reply := tg.NewMessage(msg.Chat.ID, text)

		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	post := posts[r.Intn(len(posts))-1]

	switch {
	case strings.HasSuffix(post.Image, models.FormatWebM):
		commandRandom(msg)
		return
	case strings.HasSuffix(post.Image, models.FormatGIF):
		document := tg.NewDocument(msg.Chat.ID, post.FileURL(res))
		_, err = bot.Bot.SendDocument(document)
	default:
		photo := tg.NewPhoto(msg.Chat.ID, post.FileURL(res))

		if post.Sample {
			photo.Photo = post.SampleURL(res)
		}

		_, err = bot.Bot.SendPhoto(photo)
	}
	if err != nil {
		commandRandom(msg)
	}
}
