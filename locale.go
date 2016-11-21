package main

type (
	Localization struct {
		English  Translation `json:"english"`
		Russian  Translation `json:"russian"`
		TChinese Translation `json:"tchinese"`
		SChinese Translation `json:"schinese"`
	}

	Translation struct {
		Achievements []Achievements `json:"achievements"`
		Buttons      Buttons        `json:"buttons"`
		Inline       Inline         `json:"inline"`
		Messages     Messages       `json:"messages"`
		Name         string         `json:"name"`
		Rating       Rating         `json:"rating"`
		Types        Types          `json:"types"`
	}

	Achievements struct {
		Description string `json:"description"`
		Title       string `json:"title"`
	}

	Buttons struct {
		Achievements string `json:"achievements"`
		Cancel       string `json:"cancel"`
		Channel      string `json:"channel"`
		Cheatsheet   string `json:"cheatsheet"`
		Disable      string `json:"disable"`
		Donate       string `json:"donate"`
		Enable       string `json:"enable"`
		FastStart    string `json:"fast_start"`
		Feedback     string `json:"feedback"`
		Group        string `json:"group"`
		Help         string `json:"help"`
		Language     string `json:"language"`
		More         string `json:"more"`
		Original     string `json:"original"`
		Random       string `json:"random"`
		Rate         string `json:"rate"`
		Settings     string `json:"settings"`
		Share        string `json:"share"`
	}

	Inline struct {
		NoResult InlineResult `json:"no_result"`
		Result   InlineResult `json:"result"`
	}

	InlineResult struct {
		Description string `json:"description"`
		Title       string `json:"title"`
	}

	Messages struct {
		BlushBoard string `json:"blushboard"`
		CheatSheet string `json:"cheatsheet"`
		Donate     string `json:"donate"`
		Help       string `json:"help"`
		Info       string `json:"info"`
		Start      string `json:"start"`
	}

	Rating struct {
		Explicit     string `json:"explicit"`
		Questionable string `json:"questionable"`
		Safe         string `json:"safe"`
		Unknown      string `json:"unknown"`
	}

	Types struct {
		Animation string `json:"animation"`
		Image     string `json:"image"`
		Video     string `json:"video"`
	}
)
