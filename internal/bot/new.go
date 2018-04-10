package bot

import tg "github.com/toby3d/telegram"

var Bot *tg.Bot

func New(accessToken string) (*tg.Bot, error) {
	return tg.New(accessToken)
}
