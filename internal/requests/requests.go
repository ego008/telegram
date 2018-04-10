package requests

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/HentaiDB/HentaiDBot/internal/resources"
	"github.com/HentaiDB/HentaiDBot/pkg/models"
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

var ErrNotOk = errors.New("status code is not 200")

func Results(resource string, params *Params) ([]models.Result, error) {
	cfg := resources.Resources[resource]

	requestURL := url.URL{
		Scheme: cfg.GetString("scheme"),
		Host:   cfg.GetString("host"),
		Path:   cfg.GetString("path"),
	}
	if requestURL.Scheme == "" {
		requestURL.Scheme = "http"
	}

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
