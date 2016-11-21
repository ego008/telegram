package main

import t "github.com/go-telegram-bot-api/telegram-bot-api"

type (
	// Configuration is a main config
	Configuration struct {
		DataBase DataBase   `json:"database"`
		Telegram Telegram   `json:"telegram"`
		Botan    Botan      `json:"botan"`
		Links    Links      `json:"links"`
		Resource []Resource `json:"resources"`
	}

	// Telegram API settings
	DataBase struct {
		Address  string `json:"address"`
		DataBase string `json:"database"`
		Table    string `json:"table"`
	}

	// Telegram API settings
	Telegram struct {
		Admin   int     `json:"admin"` // For future, to get feedback
		Group   int64   `json:"group"` // For easter eggs
		Token   string  `json:"token"`
		Webhook Webhook `json:"webhook"`
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

	Links struct {
		Channel string `json:"channel"`
		Donate  string `json:"donate"`
		Group   string `json:"group"`
		Rate    string `json:"rate"`
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

	MetrikaMessage struct {
		*t.Message
	}

	MetrikaInlineQuery struct {
		*t.InlineQuery
	}

	MetrikaChosenInlineResult struct {
		*t.ChosenInlineResult
	}
)
