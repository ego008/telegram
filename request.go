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

	requestURL := &url.URL{
		Scheme: resource["scheme"].(string),
		Host:   resource["host"].(string),
		Path:   resource["path"].(string),
	}

	if resource["query"] != nil {
		query := resource["query"].(map[string]interface{})
		args := requestURL.Query()

		if query["limit"] != nil &&
			req.Limit != 0 {
			args.Add(query["limit"].(string), strconv.Itoa(req.Limit))
		}
		if query["page"] != nil &&
			req.PageID > 0 {
			args.Add(query["page"].(string), strconv.Itoa(req.PageID))
		}
		if query["tags"] != nil &&
			req.Tags != "" {
			args.Add(query["tags"].(string), req.Tags)
		}
		if query["id"] != nil &&
			req.ID > 0 {
			args.Add(query["id"].(string), strconv.Itoa(req.ID))
		}
		if query["custom"] != nil {
			rawCustom := query["custom"].([]interface{})
			custom := make([]map[string]interface{}, len(rawCustom))
			for i := range rawCustom {
				custom[i] = rawCustom[i].(map[string]interface{})
			}

			for i := range custom {
				for key, val := range custom[i] {
					var value string
					switch v := val.(type) {
					case string:
						value = v
					case int:
						value = strconv.Itoa(v)
					}
					args.Add(key, value)
				}
			}
		}

		requestURL.RawQuery = args.Encode()
	}

	log.Ln("RequestURL:", requestURL.String())
	_, resp, err := http.Get(nil, requestURL.String())
	if err != nil {
		return nil, err
	}

	var posts []post
	err = json.NewDecoder().Decode(resp, &posts)
	return posts, err
}
