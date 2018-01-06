package main

import (
	"net/url"
	"strconv"

	log "github.com/kirillDanshin/dlog"
	json "github.com/pquerna/ffjson/ffjson"
	http "github.com/valyala/fasthttp"
)

type (
	params struct {
		ID, PageID, Limit int
		Tags              string
	}

	post struct {
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
		FileURL      string `json:"file_url"`
	}
)

func request(res string, req *params) ([]post, error) {
	resource := resources[res]

	var requestURL url.URL
	requestURL.Scheme = resource.UString("scheme", "http")
	requestURL.Host = resource.UString("host")
	requestURL.Path = resource.UString("path")

	args := requestURL.Query()
	args.Add("page", "dapi")
	args.Add("s", "post")
	args.Add("q", "index")
	args.Add("json", strconv.Itoa(1))

	if req.Limit != 0 {
		args.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.PageID > 0 {
		args.Add("pid", strconv.Itoa(req.PageID))
	}
	if req.Tags != "" {
		args.Add("tags", req.Tags)
	}
	if req.ID > 0 {
		args.Add("id", strconv.Itoa(req.ID))
	}

	requestURL.RawQuery = args.Encode()

	log.Ln("RequestURL:", requestURL.String())
	_, resp, err := http.Get(nil, requestURL.String())
	if err != nil {
		return nil, err
	}

	var posts []post
	err = json.NewDecoder().Decode(resp, &posts)
	return posts, err
}
