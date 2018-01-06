package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tg "github.com/toby3d/telegram"
)

func commandRandom(msg *tg.Message) {
	// usr, err := dbGetUserElseAdd(msg.From.ID, msg.From.LanguageCode)
	// errCheck(err)

	_, err := bot.SendChatAction(msg.Chat.ID, tg.ActionUploadPhoto)
	errCheck(err)

	// T, err := langSwitch(usr.Language, msg.From.LanguageCode)
	// errCheck(err)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := resourcesTags[r.Intn(len(resourcesTags)-1)]

	posts, err := request(res, &params{
		Tags: msg.CommandArgument(),
	})
	if err != nil {
		commandRandom(msg)
		return
	}

	if len(posts) <= 0 {
		text := fmt.Sprint("No results by ", msg.CommandArgument(), " tags.")
		reply := tg.NewMessage(msg.Chat.ID, text)

		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	post := posts[r.Intn(len(posts))-1]

	switch {
	case strings.HasSuffix(post.Image, "webm"):
		commandRandom(msg)
		return
	case strings.HasSuffix(post.Image, "gif"):
		document := tg.NewDocument(msg.Chat.ID, post.fileURL(res))
		_, err = bot.SendDocument(document)
	default:
		photo := tg.NewPhoto(msg.Chat.ID, post.fileURL(res))

		if post.Sample {
			photo.Photo = post.sampleURL(res)
		}

		_, err = bot.SendPhoto(photo)
	}
	if err != nil {
		commandRandom(msg)
	}
}
