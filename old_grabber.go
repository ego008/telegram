package main

/*
import (
	"net/url"
	"strconv"

	json "github.com/pquerna/ffjson/ffjson"
	http "github.com/valyala/fasthttp"
)

type (
	request struct {
		ID, PageID, ChangeID, Limit int
		Tags                        string
	}

	gPost struct {
		Change       int64  `json:"change"`
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

func getPosts(req *request) ([]gPost, error) {
	uri := url.URL{
		Scheme: "https",
		Host:   "gelbooru.com",
		Path:   "index.php",
	}

	q := uri.Query()
	q.Add("page", "dapi")
	q.Add("s", "post")
	q.Add("q", "index")
	q.Add("json", strconv.Itoa(1))
	if req.Limit > 0 {
		q.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.PageID > 0 {
		q.Add("pid", strconv.Itoa(req.PageID))
	}
	if req.ChangeID > 0 {
		q.Add("cid", strconv.Itoa(req.ChangeID))
	}
	if req.ID > 0 {
		q.Add("id", strconv.Itoa(req.ID))
	}
	if req.Tags != "" {
		q.Add("tags", req.Tags)
	}
	uri.RawQuery = q.Encode()

	_, resp, err := http.Get(nil, uri.String())
	if err != nil {
		return nil, err
	}

	var obj []gPost
	err = json.NewDecoder().Decode(resp, &obj)
	return obj, err
}
*/
