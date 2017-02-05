package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	f "github.com/valyala/fasthttp"
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

	// Post defines a structure for Gelbooru only(?)
	Post struct {
		Change       int    `json:"change"`
		Directory    string `json:"directory"`
		FileURL      string `json:"file_url"`
		Hash         string `json:"hash"`
		Height       int    `json:"height"`
		ID           int    `json:"id"`
		Image        string `json:"image"`
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
)

func getPosts(req Request) []Post {
	var args f.Args
	args.Add("page", "dapi")
	args.Add("s", "post")
	args.Add("q", "index")
	args.Add("json", strconv.Itoa(1))
	if req.Limit > 0 {
		args.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.PageID > 0 {
		args.Add("pid", strconv.Itoa(req.PageID))
	}
	if req.ChangeID > 0 {
		args.Add("cid", strconv.Itoa(req.ChangeID))
	}
	if req.ID > 0 {
		args.Add("id", strconv.Itoa(req.ID))
	}
	if req.Tags != "" {
		args.Add("tags", req.Tags)
	}

	repository := fmt.Sprintf("%s/index.php?%s", cfg["resource_url"].(string), args.String())
	_, resp, err := f.Get(nil, repository)
	if err != nil {
		log.Printf("[Bot] GET request error: %+v", err)
	}

	var obj []Post
	json.NewDecoder(bytes.NewReader(resp)).Decode(&obj)

	return obj
}
