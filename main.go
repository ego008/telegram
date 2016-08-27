package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	appMetrika = make(chan bool)
	bot        *tgbotapi.BotAPI
	config     Configuration
	metrika    botan.Botan
	resNum     = 20 // Select Gelbooru by default, remake in name search(?)
	update     tgbotapi.Update
)

func init() {
	// Read configuration
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("[Configuration] Reading error: %+v", err)
	} else {
		log.Println("[Configuration] Read successfully.")
	}
	// Decode configuration
	if err = json.Unmarshal(file, &config); err != nil {
		log.Fatalf("[Configuration] Decoding error: %+v", err)
	}

	// Initialize bot
	newBot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Panicf("[Bot] Initialize error: %+v", err)
	} else {
		bot = newBot
		log.Printf("[Bot] Authorized as @%s", bot.Self.UserName)
	}

	metrika = botan.New(config.Botan.Token) // Initialize botan
	log.Println("[Botan] ACTIVATED")
}

func main() {
	startUptime := time.Now() // Set start UpTime time

	debugMode := flag.Bool("debug", false, "enable debug logs")
	webhookMode := flag.Bool("webhook", false, "enable webhooks support")
	flag.Parse()

	bot.Debug = *debugMode

	updates := make(<-chan tgbotapi.Update)
	if *webhookMode == true {
		if _, err := bot.SetWebhook(tgbotapi.NewWebhook(config.Telegram.Webhook.Set + config.Telegram.Token)); err != nil {
			log.Printf("Set webhook error: %+v", err)
		}
		updates = bot.ListenForWebhook(config.Telegram.Webhook.Listen + config.Telegram.Token)
		go http.ListenAndServe(config.Telegram.Webhook.Serve, nil)
	} else {
		upd := tgbotapi.NewUpdate(0)
		upd.Timeout = 60
		updater, err := bot.GetUpdatesChan(upd)
		if err != nil {
			log.Printf("[Bot] Getting updates error: %+v", err)
		}
		updates = updater
	}

	// Updater
	for update = range updates {
		log.Printf("[Bot] Update response: %+v", updates)

		// Chat actions
		if update.Message != nil {
			switch true {
			case update.Message.Command() == "start" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)): // Requirement Telegram platform
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/start", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /start %s", answer.Status)
					appMetrika <- true
				})

				bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)) // Force feedback

				message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(startMessage, update.Message.From.FirstName))
				message.ParseMode = "markdown"
				message.DisableWebPagePreview = true
				message.ReplyToMessageID = update.Message.MessageID
				if update.Message.Chat.IsPrivate() == true {
					message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							// Showing tutorial button only in private chat for demonstration work
							tgbotapi.NewInlineKeyboardButtonSwitch("See how to do this!", "hatsune_miku rating:safe"),
						),
					)
				}
				if _, err := bot.Send(message); err != nil {
					log.Printf("[Bot] Sending message error: %+v", err)
				}

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "help" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)): // Requirement Telegram platform
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/help", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /help %s", answer.Status)
					appMetrika <- true
				})

				bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)) // Force feedback

				message := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
				message.ParseMode = "markdown"
				message.DisableWebPagePreview = true
				message.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(message); err != nil {
					log.Printf("[Bot] Sending message error: %+v", err)
				}

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "cheatsheet" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)):
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/cheatsheet", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /cheatsheet %s", answer.Status)
					appMetrika <- true
				})

				bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)) // Force feedback

				// For now - get Cheat Sheet from Gelbooru
				message := tgbotapi.NewMessage(update.Message.Chat.ID, cheatSheetMessage)
				message.ParseMode = "markdown"
				message.DisableWebPagePreview = true
				message.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(message); err != nil {
					log.Printf("[Bot] Sending message error: %+v", err)
				}

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "random" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)):
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/random", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /random %s", answer.Status)
					appMetrika <- true
				})

				bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatUploadDocument)) // Force feedback

				randomSource := rand.NewSource(time.Now().UnixNano()) // Maximum randomizing dice
				totalPosts := getPosts(Request{ID: 0})                // Get last upload post
				random := rand.New(randomSource)                      // Create magical dice
				var randomFile []Post

				for {
					randomPost := random.Intn(totalPosts[0].ID)    // Generate a random ID number from first to last ID post
					randomFile = getPosts(Request{ID: randomPost}) // Call to selected ID
					if len(randomFile) > 0 {
						break // If post is NOT blocked or erroneous
					}
					log.Println("[Bot] This is not image. Reroll dice!")
				}

				var button tgbotapi.InlineKeyboardMarkup
				if update.Message.Chat.IsPrivate() == true {
					button = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonURL("Original image", randomFile[0].FileURL),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonSwitch("Share", "id:"+strconv.Itoa(randomFile[0].ID)),
						),
					)
				} else {
					button = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonURL("Original image", randomFile[0].FileURL),
						),
					)
				}

				_, body, err := fasthttp.Get(nil, randomFile[0].FileURL)
				if err != nil {
					log.Printf("[Bot] Get random image by URL error: %+v", err)
				}
				bytes := tgbotapi.FileBytes{Name: randomFile[0].Image, Bytes: body}

				switch {
				case strings.Contains(randomFile[0].FileURL, ".mp4"):
					bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatUploadVideo)) // Force feedback

					message := tgbotapi.NewVideoUpload(update.Message.Chat.ID, bytes)
					message.ReplyToMessageID = update.Message.MessageID
					message.ReplyMarkup = &button
					if _, err := bot.Send(message); err != nil {
						log.Printf("[Bot] Sending message error: %+v", err)
					}
				case strings.Contains(randomFile[0].FileURL, ".gif") || strings.Contains(randomFile[0].FileURL, ".webm"):
					bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatUploadDocument)) // Force feedback

					message := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, bytes)
					message.ReplyToMessageID = update.Message.MessageID
					message.ReplyMarkup = &button
					if _, err := bot.Send(message); err != nil {
						log.Printf("[Bot] Sending message error: %+v", err)
					}
				default:
					bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatUploadPhoto)) // Force feedback

					message := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, bytes)
					message.ReplyToMessageID = update.Message.MessageID
					message.ReplyMarkup = &button
					if _, err := bot.Send(message); err != nil {
						log.Printf("[Bot] Sending message error: %+v", err)
					}
				}

				<-appMetrika // Send track to Yandex.AppMetrika
			case update.Message.Command() == "info" && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)):
				// Track action
				metrika.TrackAsync(update.Message.From.ID, MetrikaMessage{update.Message}, "/info", func(answer botan.Answer, err []error) {
					log.Printf("[Botan] Track /info %s", answer.Status)
					appMetrika <- true
				})

				bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatUploadPhoto)) // Force feedback

				uptime := time.Since(startUptime).String()

				// For now - get Cheat Sheet from Gelbooru
				message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(infoMessage, uptime))
				message.ParseMode = "markdown"
				message.DisableWebPagePreview = true
				message.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(message); err != nil {
					log.Printf("[Bot] Sending message error: %+v", err)
				}

				<-appMetrika // Send track to Yandex.AppMetrika
			default:
				getEasterEgg() // Secret actions and commands ;)
			}
		}

		// Inline actions
		if update.InlineQuery != nil {
			// Track action
			metrika.TrackAsync(update.InlineQuery.From.ID, MetrikaInlineQuery{update.InlineQuery}, "Search", func(answer botan.Answer, err []error) {
				log.Printf("[Botan] Track Search %s", answer.Status)
				appMetrika <- true
			})

			// Check result pages
			var posts []Post
			var resultPage int
			if len(update.InlineQuery.Offset) > 0 {
				resultPage, _ = strconv.Atoi(update.InlineQuery.Offset)
				posts = getPosts(Request{PageID: resultPage, Tags: update.InlineQuery.Query})
			} else {
				posts = getPosts(Request{Tags: update.InlineQuery.Query})
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

					button := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonURL("Original image", posts[i].FileURL),
						),
					)

					switch {
					case strings.Contains(posts[i].FileURL, ".webm"): // It is necessary to get around error 403 when requesting video :|
						// query := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(i), posts[i].FileURL) // Does not work
						// query.MimeType = "text/html" // Link on widget-page?
						// query.MimeType = "video/mp4" // Does not work for .webm
						// query.ThumbURL = preview
						// query.Width = posts[i].Width
						// query.Height = posts[i].Height
						// query.Title = "Video by " + strings.Title(posts[i].Owner)
						// query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
						// query.ReplyMarkup = &button
						// result = append(result, query)
						continue
					case strings.Contains(posts[i].FileURL, ".mp4"): // Just in case. Why not? ¯\_(ツ)_/¯
						query := tgbotapi.NewInlineQueryResultVideo(strconv.Itoa(i), posts[i].FileURL)
						query.MimeType = "video/mp4"
						query.ThumbURL = preview
						query.Width = posts[i].Width
						query.Height = posts[i].Height
						query.Title = "Video by " + strings.Title(posts[i].Owner)
						query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
						query.ReplyMarkup = &button
						result = append(result, query)
					case strings.Contains(posts[i].FileURL, ".gif"):
						query := tgbotapi.NewInlineQueryResultGIF(strconv.Itoa(i), posts[i].FileURL)
						query.ThumbURL = posts[i].FileURL
						query.Width = posts[i].Width
						query.Height = posts[i].Height
						query.Title = "Animation by " + strings.Title(posts[i].Owner)
						query.ReplyMarkup = &button
						result = append(result, query)
					default:
						query := tgbotapi.NewInlineQueryResultPhoto(strconv.Itoa(i), posts[i].FileURL)
						query.ThumbURL = preview
						query.Width = posts[i].Width
						query.Height = posts[i].Height
						query.Title = "Image by " + strings.Title(posts[i].Owner)
						query.Description = "Rating: " + rating + "\nScore: " + strconv.Itoa(posts[i].Score) + "\nTags: " + posts[i].Tags
						query.ReplyMarkup = &button
						result = append(result, query)
					}
				}
			case len(posts) == 0: // Found nothing
				query := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, noInlineResultTitle, noInlineResultMessage)
				query.Description = noInlineResultDescription
				result = append(result, query)
			}

			// Configure inline-mode
			inlineConfig := tgbotapi.InlineConfig{}
			inlineConfig.InlineQueryID = update.InlineQuery.ID
			inlineConfig.IsPersonal = true
			inlineConfig.CacheTime = 0
			inlineConfig.Results = result
			if len(posts) == 50 {
				inlineConfig.NextOffset = strconv.Itoa(resultPage + 1) // If available next page of results
			}

			if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
				log.Printf("[Bot] Answer inline-query error: %+v", err)
			}

			<-appMetrika // Send track to Yandex.AppMetrika
		}

		if update.ChosenInlineResult != nil {
			metrika.TrackAsync(update.ChosenInlineResult.From.ID, MetrikaChosenInlineResult{update.ChosenInlineResult}, "Find", func(answer botan.Answer, err []error) {
				log.Printf("[Botan] Track Find %s", answer.Status)
				appMetrika <- true
			})
			<-appMetrika // Send track to Yandex.AppMetrika
		}

		if update.CallbackQuery != nil {
			metrika.TrackAsync(update.ChosenInlineResult.From.ID, MetrikaCallbackQuery{update.CallbackQuery}, "Action", func(answer botan.Answer, err []error) {
				log.Printf("[Botan] Track Find %s", answer.Status)
				appMetrika <- true
			})
			<-appMetrika // Send track to Yandex.AppMetrika
		}
	}
}
