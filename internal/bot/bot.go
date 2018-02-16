package bot

import (
	"github.com/HentaiDB/HentaiDBot/internal/config"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	tg "github.com/toby3d/telegram"
)

var Bot *tg.Bot

func Initialize() {
	var err error
	Bot, err = tg.NewBot(config.Config.UString("telegram.token"))
	errors.Check(err)
}
