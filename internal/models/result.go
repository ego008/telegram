package models

type Result struct {
	Source       string `json:"source"`
	Directory    string `json:"directory"`
	Hash         string `json:"hash"`
	Height       int    `json:"height"`
	ID           int    `json:"id"`
	Image        string `json:"image"`
	Change       int64  `json:"change"`
	Owner        string `json:"owner"`
	ParentID     int    `json:"parent_id"`
	Rating       string `json:"rating"`
	Sample       bool   `json:"sample"`
	SampleHeight int    `json:"sample_height"`
	SampleWidth  int    `json:"sample_width"`
	Score        int    `json:"score"`
	Tags         string `json:"tags"`
	Width        int    `json:"width"`
}
