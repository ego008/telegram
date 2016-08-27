package main

type Localization struct {
	English Translation `json:"english"`
	Russian Translation `json:"russian"`
}

type Translation struct {
	Achievements []Achievements `json:"achievements"`
	Buttons      Buttons        `json:"buttons"`
	Inline       Inline         `json:"inline"`
	Messages     Messages       `json:"messages"`
	Name         string         `json:"name"`
	Rating       Rating         `json:"rating"`
	Types        Types          `json:"types"`
}

type Achievements struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Buttons struct {
	Achievements string `json:"achievements"`
	Cancel       string `json:"cancel"`
	Cheatsheet   string `json:"cheatsheet"`
	Disable      string `json:"disable"`
	Enable       string `json:"enable"`
	Feedback     string `json:"feedback"`
	GIF          string `json:"gif"`
	Help         string `json:"help"`
	How          string `json:"how"`
	Image        string `json:"image"`
	Language     string `json:"language"`
	More         string `json:"more"`
	Random       string `json:"random"`
	Settings     string `json:"settings"`
	Share        string `json:"share"`
	Video        string `json:"video"`
}

type Inline struct {
	NoResult          NoInlineResult `json:"no_result"`
	ResultDescription string         `json:"result_description"`
}

type NoInlineResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

type Messages struct {
	Info       string `json:"info"`
	CheatSheet string `json:"cheatsheet"`
	Help       string `json:"help"`
	Start      string `json:"start"`
}

type Rating struct {
	Explicit     string `json:"explicit"`
	Questionable string `json:"questionable"`
	Safe         string `json:"safe"`
	Unknown      string `json:"unknown"`
}

type Types struct {
	Image     string `json:"image"`
	Animation string `json:"animation"`
	Video     string `json:"video"`
}
