package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

// Arguments for getPosts()
type Request struct {
	Limit    int
	PageID   int
	Tags     string
	ChangeID int
	ID       int
}

// Configuration is a main config
type Configuration struct {
	Telegram Telegram   `json:"telegram"`
	Botan    Botan      `json:"botan"`
	Resource []Resource `json:"resources"`
}

// Telegram API settings
type Telegram struct {
	Webhook     Webhook `json:"webhook"`
	Token      string `json:"token"`
	Admin      int64  `json:"admin"`      // For future, to get feedback
	SuperGroup int64  `json:"supergroup"` // For easter eggs
	Channel    string `json:"channel"`
}

type Webhook struct {
	Set    string `json:"set"`
	Listen string `json:"listen"`
	Serve  string `json:"serve"`
}

// Botan structure defines botan API settings
type Botan struct {
	Token string `json:"token"`
}

// Resource structure
type Resource struct {
	Name     string   `json:"name"`
	Settings Settings `json:"settings"`
}

// Settings structure defines resource settings
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

// Post defines a structure for Danbooru only(?)
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
	Score        int    `json:"score"`
	Tags         string `json:"tags"`
	Width        int    `json:"width"`
	FileURL      string `json:"file_url"`
}

type MetrikaMessage struct {
	*tgbotapi.Message
}

type MetrikaInlineQuery struct {
	*tgbotapi.InlineQuery
}

type MetrikaChosenInlineResult struct {
	*tgbotapi.ChosenInlineResult
}
