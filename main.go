package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	config Configuration
	bot    *tgbotapi.BotAPI
	resNum = 20 // Select Gelbooru by default, remake in name search(?)
)

func init() {
	// Read configuration
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panicf("Error reading configuration file: %s", err)
	} else {
		log.Println("Сonfiguration file is read successfully.")
	}
	// Decode configuration
	if err = json.Unmarshal(file, &config); err != nil {
		log.Panicf("Error decoding configuration file: %s", err)
	}

	// Initialize bot
	newBot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Panic(err)
	} else {
		bot = newBot
		bot.Debug = true
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}
}

func main() {
	// Yandex.AppMetrika
	appMetrika := make(chan bool)
	botanio := botan.New(config.Botan.Token)

	// Timer updates (webhooks works only in production)
	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	updates, err := bot.GetUpdatesChan(upd)
	if err != nil {
		log.Fatalf("Error getting updates: %s", err)
	}

	// Updater
	for update := range updates {
		log.Printf("%+v", update)

		// Chat actions
		if update.Message != nil {
			switch update.Message.Text {
			case "/start": // Requirement Telegram platform
				// Track action
				botanio.TrackAsync(update.Message.From.ID, update, "/start", func(answer botan.Answer, err []error) {
					log.Printf("Asynchonous: %+v", answer)
					appMetrika <- true
				})

				message := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					fmt.Sprintf(
						startMsgTpl,
						update.Message.From.FirstName,
					),
				)
				message.DisableWebPagePreview = true
				message.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(message); err != nil {
					log.Fatal(err)
				}

				// Send track to Yandex.AppMetrika
				<-appMetrika
			case "/help": // Requirement Telegram platform
				// Track action
				botanio.TrackAsync(update.Message.From.ID, update, "/help", func(answer botan.Answer, err []error) {
					log.Printf("Asynchonous: %+v", answer)
					appMetrika <- true
				})

				// For now - get Cheat Sheet from Gelbooru
				// It will be transferred to command like /cheatsheet
				message := tgbotapi.NewMessage(update.Message.Chat.ID, helpMsg)
				message.ParseMode = "html"
				message.DisableWebPagePreview = true
				message.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(message); err != nil {
					log.Fatal(err)
				}

				// Send track to Yandex.AppMetrika
				<-appMetrika
			default:
				GetEasterEgg(bot, botanio, update) // Secret actions and commands
			}
		}

		// Inline actions
		if update.InlineQuery != nil {
			// Track action
			// It is necessary to fix <nil> tracking ChosenInlineResult. :\
			botanio.TrackAsync(update.InlineQuery.From.ID, update, "inline", func(answer botan.Answer, err []error) {
				log.Printf("Asynchonous: %+v", answer)
				appMetrika <- true
			})

			// Check result pages
			var posts []Post
			var resultPage int
			if len(update.InlineQuery.Offset) > 0 {
				posts = getPosts(update.InlineQuery.Query, update.InlineQuery.Offset)
				resultPage, _ = strconv.Atoi(update.InlineQuery.Offset)
			} else {
				posts = getPosts(update.InlineQuery.Query, "")
			}

			// Analysis of results
			var result []interface{}
			switch {
			case len(posts) > 0:
				for i := 0; i < len(posts); i++ {
					// Universal(?) preview url
					preview := config.Resource[resNum].Settings.URL + config.Resource[resNum].Settings.ThumbsDir + posts[i].Directory + config.Resource[resNum].Settings.ThumbsPart + posts[i].Hash + ".jpg"

					// Rating
					var rating string
					switch posts[i].Rating {
					case "s":
						rating = "Safe"
					case "e":
						rating = "Explicit"
					case "q":
						rating = "Questionable"
					default:
						rating = "Unknown"
					}

					// URL-button with a direct link to result
					button := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Original image", posts[i].FileURL)))

					switch {
					case strings.Contains(posts[i].FileURL, ".webm"): // It is necessary to get around error 403 when requesting video :|
						// query := tgbotapi.NewInlineQueryResultVideo(update.InlineQuery.ID+strconv.Itoa(posts[i].ID), posts[i].FileURL) // Does not work
						// query.MimeType = "text/html" // Link on widget-page?
						// query.MimeType = "video/mp4" // Does not work for .webm
						// query.ThumbURL = preview
						// query.Width = posts[i].Width
						// query.Height = posts[i].Height
						// query.Title = "by " + strings.Title(posts[i].Owner)
						// query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
						// query.ReplyMarkup = &button
						// result = append(result, query)
						continue
					case strings.Contains(posts[i].FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
						query := tgbotapi.NewInlineQueryResultVideo(update.InlineQuery.ID+strconv.Itoa(posts[i].ID), posts[i].FileURL)
						query.MimeType = "video/mp4"
						query.ThumbURL = preview
						query.Width = posts[i].Width
						query.Height = posts[i].Height
						query.Title = "by " + strings.Title(posts[i].Owner)
						query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
						query.ReplyMarkup = &button
						result = append(result, query)
					case strings.Contains(posts[i].FileURL, ".gif"):
						query := tgbotapi.NewInlineQueryResultGIF(update.InlineQuery.ID+strconv.Itoa(posts[i].ID), posts[i].FileURL)
						query.ThumbURL = posts[i].FileURL
						query.Width = posts[i].Width
						query.Height = posts[i].Height
						query.Title = "by " + strings.Title(posts[i].Owner)
						query.ReplyMarkup = &button
						result = append(result, query)
					default:
						query := tgbotapi.NewInlineQueryResultPhoto(update.InlineQuery.ID+strconv.Itoa(posts[i].ID), posts[i].FileURL)
						query.ThumbURL = preview
						query.Width = posts[i].Width
						query.Height = posts[i].Height
						query.Title = "by " + strings.Title(posts[i].Owner)
						query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
						query.ReplyMarkup = &button
						result = append(result, query)
					}
				}
			case len(posts) == 0: // Found nothing
				query := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Nobody here but us chickens!", "Sumimasen, but, unfortunately I could not find desired content. 😓\nBut perhaps this it already present in @HentaiDB channel.")
				query.Description = "Try search a different combination of tags."
				result = append(result, query)
			}

			// Configure inline-mode
			inlineConfig := tgbotapi.InlineConfig{}
			inlineConfig.InlineQueryID = update.InlineQuery.ID
			inlineConfig.IsPersonal = true
			inlineConfig.CacheTime = 0
			inlineConfig.Results = result
			// If available next page of results
			if len(posts) == 50 {
				inlineConfig.NextOffset = strconv.Itoa(resultPage + 1)
			}

			if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
				log.Fatal(err)
			}

			<-appMetrika // Send track to Yandex.AppMetrika
		}
	}
}
