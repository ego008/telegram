package updates

import (
	"github.com/HentaiDB/HentaiDBot/internal/callbacks"
	"github.com/HentaiDB/HentaiDBot/internal/database"
	"github.com/HentaiDB/HentaiDBot/internal/inline"
	"github.com/HentaiDB/HentaiDBot/internal/messages"
)

func Updates() {
	defer database.Close()

	updatesChannel := getUpdatesChannel()
	for update := range updatesChannel {
		switch {
		case update.Message != nil:
			messages.Message(update.Message)
		case update.InlineQuery != nil &&
			len(update.InlineQuery.Query) <= 255:
			// inline.InlineQuery(update.InlineQuery)
		case update.ChosenInlineResult != nil:
			// ChosenInlineResult(update.ChosenInlineResult)
		case update.CallbackQuery != nil:
			// callbacks.CallbackQuery(update.CallbackQuery)
		case update.ChannelPost != nil:
			// channelPost(update.ChannelPost)
		}
	}
}
