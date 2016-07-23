package main

import (
	"bytes"
	"encoding/json"
	"github.com/kirillDanshin/myutils"
	"github.com/valyala/fasthttp"
	"log"
)

// Universal(?) function obtain content
func getPosts(tags string, pid string) []Post {
	repository := myutils.Concat(config.Resource[resNum].Settings.URL, "index.php?page=dapi&s=post&q=index&json=1&limit=50") // JSON API with 50 results (Telegram limit)
	if tags != "" {
		repository = myutils.Concat(repository, "&tags=", tags) // Insert tags
	}
	if pid != "" {
		repository = myutils.Concat(repository, "&pid=", pid) // Insert result-page
	}
	_, resp, err := fasthttp.Get(nil, repository)
	if err != nil {
		log.Printf("[Bot] GET request error: %+v", err)
	}
	var obj []Post
	json.NewDecoder(bytes.NewReader(resp)).Decode(&obj)
	return obj
}
