package main

import (
	"encoding/json"
	"github.com/botanio/sdk/go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	//"github.com/valyala/fasthttp" // Need to replace "net/http" later
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	config Configuration
	bot    *tgbotapi.BotAPI
	resNum = 20 // Select Gelbooru by default, remake in name search(?)
)

type Configuration struct {
	Telegram Telegram   `json:"telegram"`
	Botan    Botan      `json:"botan"`
	Resource []Resource `json:"resources"`
}

type Telegram struct {
	Token       string `json:"token"`
	AdminChatID int    `json:"admin_chat_id"` // For future, to get feedback
}

type Botan struct {
	Token string `json:"token"`
}

type Resource struct {
	Name     string   `json:"name"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	URL        string `json:"url"`
	Template   string `json:"template,omniempty"`   // For future(?)
	CheatSheet string `json:"cheatsheet,omniempty"` // For future, for parce help instructions
	ThumbsDir  string `json:"thumbs_dir,omniempty"`
	ImagesDir  string `json:"images_dir,omniempty"`
	ThumbsPart string `json:"thumbs_part,omniempty"`
	ImagesPart string `json:"images_part,omniempty"`
	AddPath    string `json:"addpath,omniempty"` // ???
}

// Structure for Danbooru only(?)
type Post struct {
	Directory    string `json:"directory"`
	Hash         string `json:"hash"`
	Height       int    `json:"height"`
	ID           int    `json:"id"`
	Image        string `json:"image"`
	Change       int    `json:"change"`
	Owner        string `json:"owner"`
	ParentID     int    `json:"parent_id"`
	Rating       string `json:"rating"`
	Sample       string `json:"sample"`
	SampleHeight int    `json:"sample_height"`
	SampleWidth  int    `json:"sample_width"`
	Score        int    `json:"score""`
	Tags         string `json:"tags"`
	Width        int    `json:"width"`
	FileURL      string `json:"file_url"`
}

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

				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi "+update.Message.From.FirstName+"!\n\nThis is the official @HentaiDB bot. You can browse the Danbooru pics, GIF's and videos here.\n\nYou can also use it to search and share content with your friends. Just type \"@HentaiDBot hatsune_miku\" in any chat and select the result you want to send.")
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
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "<b>tag1 tag2</b>\n<code>Search for posts that have tag1 and tag2.</code>\n\n<b>~tag1 ~tag2</b>\n<code>Search for posts that have tag1 or tag2. (Currently does not work)</code>\n\n<b>night~</b>\n<code>Fuzzy search for the tag night. This will return results such as night fight bright and so on according to the </code><a href=\"https://en.wikipedia.org/wiki/Levenshtein_distance\">Levenshtein distance</a><code>.</code>\n\n<b>-tag1</b>\n<code>Search for posts that don't have tag1.</code>\n\n<b>ta*1</b>\n<code>Search for posts with tags that starts with ta and ends with 1.</code>\n\n<b>user:bob</b>\n<code>Search for posts uploaded by the user Bob.</code>\n\n<b>md5:foo</b>\n<code>Search for posts with the MD5 hash foo.</code>\n\n<b>md5:foo*</b>\n<code>Search for posts whose MD5 starts with the MD5 hash foo.</code>\n\n<b>rating:questionable</b>\n<code>Search for posts that are rated questionable.</code>\n\n<b>-rating:questionable</b>\n<code>Search for posts that are not rated questionable.</code>\n\n<b>parent:1234</b>\n<code>Search for posts that have 1234 as a parent (and include post 1234).</code>\n\n<b>rating:questionable rating:safe</b>\n<code>In general, combining the same metatags (the ones that have colons in them) will not work.</code>\n\n<b>rating:questionable parent:100</b>\n<code>You can combine different metatags, however.</code>\n\n<b>width:&gt;=1000 height:&gt;1000</b>\n<code>Find images with a width greater than or equal to 1000 and a height greater than 1000.</code>\n\n<b>score:&gt;=10</b>\n<code>Find images with a score greater than or equal to 10. This value is updated once daily at 12AM CST.</code>\n\n<b>sort:updated:desc</b>\n<code>Sort posts by their most recently updated order.</code>\n\n<b>Other sortable types:</b>\n<code>- id\n- score\n- rating\n- user\n- height\n- width\n- parent\n- source\n- updated\nCan be sorted by both asc or desc.</code>")
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
			var resultPage int = 0
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

// Universal(?) function obtain content
func getPosts(tags string, pid string) []Post {
	// JSON API with 50 results (Telegram limit)
	repository := config.Resource[resNum].Settings.URL + "index.php?page=dapi&s=post&q=index&json=1&limit=50"
	if tags != "" {
		repository += "&tags=" + tags // Insert tags
	}
	if pid != "" {
		repository += "&pid=" + pid // Insert result-page
	}
	resp, err := http.Get(repository) // Need to replace on "fasthttp" later :\
	if err != nil {
		log.Fatal("Error in GET request: %s", err)
	}
	defer resp.Body.Close()
	var obj []Post
	json.NewDecoder(resp.Body).Decode(&obj)
	return obj
}
