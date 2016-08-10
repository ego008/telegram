package main

import (
	"bytes"
	"encoding/json"
	"github.com/kirillDanshin/myutils"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

// Universal(?) function obtain content
func getPosts(req Request) []Post {
	repository := myutils.Concat(config.Resource[resNum].Settings.URL, "index.php?page=dapi&s=post&q=index&json=1") // JSON API with 50 results (Telegram limit)
	if req.Limit == 0 {
		req.Limit = 50
	}
	switch {
	case req.Limit > 0 && req.Limit <= 50:
		repository = myutils.Concat(repository, "&limit=", strconv.Itoa(req.Limit))
		fallthrough
	case req.PageID > 0:
		repository = myutils.Concat(repository, "&pid=", strconv.Itoa(req.PageID))
		fallthrough
	case req.Tags != "":
		repository = myutils.Concat(repository, "&tags=", req.Tags)
		fallthrough
	case req.ChangeID > 0:
		repository = myutils.Concat(repository, "&cid=", strconv.Itoa(req.ChangeID))
		fallthrough
	case req.ID > 0:
		repository = myutils.Concat(repository, "&id=", strconv.Itoa(req.ID))
	}
	_, resp, err := fasthttp.Get(nil, repository)
	if err != nil {
		log.Printf("[Bot] GET request error: %+v", err)
	}
	var obj []Post
	json.NewDecoder(bytes.NewReader(resp)).Decode(&obj)
	return obj
}
