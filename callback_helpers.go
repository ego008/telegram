package main

import tg "github.com/toby3d/go-telegram"

func callbackAlert(call *tg.CallbackQuery, text string) {
	answer := tg.NewAnswerCallbackQuery(call.ID)
	answer.Text = text

	_, err := bot.AnswerCallbackQuery(answer)
	errCheck(err)
}
