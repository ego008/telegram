package main

import (
	"bytes"
	"encoding/json"
	"github.com/kirillDanshin/myutils"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"strconv"
)

type (
	// Arguments for getPosts()
	Request struct {
		Limit    int
		PageID   int
		Tags     string
		ChangeID int
		ID       int
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
)

// Universal(?) function obtain content
func getPosts(req Request) []Post {
	repository := myutils.Concat(config.Resource[20].Settings.URL, "index.php?page=dapi&s=post&q=index&json=1")
	switch {
	case req.Limit > 0:
		repository = myutils.Concat(repository, "&limit=", strconv.Itoa(req.Limit))
		fallthrough
	case req.PageID > 0:
		repository = myutils.Concat(repository, "&pid=", strconv.Itoa(req.PageID))
		fallthrough
	case req.Tags != "":
		repository = myutils.Concat(repository, "&tags=", url.QueryEscape(req.Tags))
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
