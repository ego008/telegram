package callbacks

import (
	"github.com/HentaiDB/HentaiDBot/internal/bot"
	"github.com/HentaiDB/HentaiDBot/internal/errors"
	tg "github.com/toby3d/telegram"
)

func CallbackAlert(call *tg.CallbackQuery, text string) {
	answer := tg.NewAnswerCallbackQuery(call.ID)
	answer.Text = text

	_, err := bot.Bot.AnswerCallbackQuery(answer)
	errors.Check(err)
}
