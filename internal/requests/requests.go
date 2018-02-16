package requests

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/HentaiDB/HentaiDBot/internal/models"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	log "github.com/kirillDanshin/dlog"
	json "github.com/pquerna/ffjson/ffjson"
	http "github.com/valyala/fasthttp"
)

type Params struct {
	ID,
	PageID,
	Limit int
	Tags string
}

var ErrNotOk = errors.New("Status code is not 200")

func Results(resource string, params *Params) ([]models.Result, error) {
	res := resources.Resources[resource]

	var requestURL url.URL
	requestURL.Scheme = res.UString("scheme", "http")
	requestURL.Host = res.UString("host")
	requestURL.Path = res.UString("path")

	args := requestURL.Query()
	args.Add("page", "dapi")
	args.Add("s", "post")
	args.Add("q", "index")
	args.Add("json", strconv.Itoa(1))

	if params.Limit != 0 {
		args.Add("limit", strconv.Itoa(params.Limit))
	}
	if params.PageID > 0 {
		args.Add("pid", strconv.Itoa(params.PageID))
	}
	if params.Tags != "" {
		args.Add("tags", params.Tags)
	}
	if params.ID > 0 {
		args.Add("id", strconv.Itoa(params.ID))
	}

	requestURL.RawQuery = args.Encode()

	log.Ln("RequestURL:", requestURL.String())
	code, resp, err := http.Get(nil, requestURL.String())
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, ErrNotOk
	}

	var results []models.Result
	err = json.Unmarshal(resp, &results)
	return results, err
}
