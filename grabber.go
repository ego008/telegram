package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	// "net/url"
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
	var args f.Args
	args.Add("page", "dapi")
	args.Add("s", "post")
	args.Add("q", "index")
	args.Add("json", strconv.Itoa(1))
	switch {
	case req.Limit > 0:
		args.Add("limit", strconv.Itoa(req.Limit))
		fallthrough
	case req.PageID > 0:
		args.Add("pid", strconv.Itoa(req.PageID))
		fallthrough
	case req.ChangeID > 0:
		args.Add("cid", strconv.Itoa(req.ChangeID))
		fallthrough
	case req.ID > 0:
		args.Add("id", strconv.Itoa(req.ID))
		fallthrough
	case req.Tags != "":
		args.Add("tags", req.Tags)
	}
	repository := fmt.Sprintf("%s/index.php?%s", cfg["resource_url"].(string), args.String())
	log.Println(repository)
	_, resp, err := f.Get(nil, repository)
	if err != nil {
		log.Printf("[Bot] GET request error: %+v", err)
	}
	var obj []Post
	json.NewDecoder(bytes.NewReader(resp)).Decode(&obj)
	return obj
}
