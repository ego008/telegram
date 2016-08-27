package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

type (
	// Arguments for getPosts()
	Request struct {
		Limit    int
		PageID   int
		Tags     string
		ChangeID int
		ID       int
	}

	// Configuration is a main config
	Configuration struct {
		Telegram Telegram   `json:"telegram"`
		Botan    Botan      `json:"botan"`
		Resource []Resource `json:"resources"`
	}

	// Telegram API settings
	Telegram struct {
		Token      string  `json:"token"`
		Admin      int64   `json:"admin"`      // For future, to get feedback
		SuperGroup int64   `json:"supergroup"` // For easter eggs
		Channel    string  `json:"channel"`
		Webhook    Webhook `json:"webhook"`
	}

	Webhook struct {
		Set    string `json:"set"`
		Listen string `json:"listen"`
		Serve  string `json:"serve"`
	}

	// Botan structure defines botan API settings
	Botan struct {
		Token string `json:"token"`
	}

	// Resource structure
	Resource struct {
		Name     string   `json:"name"`
		Settings Settings `json:"settings"`
	}

	// Settings structure defines resource settings
	Settings struct {
		URL        string `json:"url"`
		Template   string `json:"template,omniempty"`   // For future(?)
		CheatSheet string `json:"cheatsheet,omniempty"` // For future, for parce help instructions
		ThumbsDir  string `json:"thumbs_dir,omniempty"`
		ImagesDir  string `json:"images_dir,omniempty"`
		ThumbsPart string `json:"thumbs_part,omniempty"`
		ImagesPart string `json:"images_part,omniempty"`
		AddPath    string `json:"addpath,omniempty"` // ???
	}

	// Post defines a structure for Danbooru only(?)
	Post struct {
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
		Score        int    `json:"score"`
		Tags         string `json:"tags"`
		Width        int    `json:"width"`
		FileURL      string `json:"file_url"`
	}

	MetrikaMessage struct {
		*tgbotapi.Message
	}

	MetrikaInlineQuery struct {
		*tgbotapi.InlineQuery
	}

	MetrikaChosenInlineResult struct {
		*tgbotapi.ChosenInlineResult
	}

	MetrikaCallbackQuery struct {
		*tgbotapi.CallbackQuery
	}
)
