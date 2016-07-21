package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/kirillDanshin/myutils"
	"github.com/valyala/fasthttp"
)

// Universal(?) function obtain content
func getPosts(tags string, pid string) []Post {
	// JSON API with 50 results (Telegram limit)
	repository := myutils.Concat(config.Resource[resNum].Settings.URL, "index.php?page=dapi&s=post&q=index&json=1&limit=50")
	if tags != "" {
		repository = myutils.Concat(repository, "&tags=", tags) // Insert tags
	}
	if pid != "" {
		repository = myutils.Concat(repository, "&pid=", pid) // Insert result-page
	}
	_, resp, err := fasthttp.Get(nil, repository)
	if err != nil {
		log.Printf("Error in GET request: %s", err)
	}
	var obj []Post
	json.NewDecoder(bytes.NewReader(resp)).Decode(&obj)
	return obj
}
