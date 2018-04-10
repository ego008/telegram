package models

import (
	"fmt"
	"net/url"
	"strings"
)

func (usr *User) GetRatingsFilter() string {
	switch {
	case usr.Ratings.Safe && !usr.Ratings.Questionable && !usr.Ratings.Exlplicit:
		// safe only
		return "rating:safe"

	case !usr.Ratings.Safe && usr.Ratings.Questionable && !usr.Ratings.Exlplicit:
		// questionable only
		return "rating:questionable"

	case !usr.Ratings.Safe && !usr.Ratings.Questionable && usr.Ratings.Exlplicit:
		// explicit only
		return "rating:explicit"

	case usr.Ratings.Safe && usr.Ratings.Questionable && !usr.Ratings.Exlplicit:
		// safe + questionable
		return "-rating:explicit"

	case usr.Ratings.Safe && !usr.Ratings.Questionable && usr.Ratings.Exlplicit:
		// safe + explicit
		return "-rating:questionable"

	case !usr.Ratings.Safe && usr.Ratings.Questionable && usr.Ratings.Exlplicit:
		// questionable + explicit
		return "-rating:safe"

	default:
		// all ratings enabled/diabled
		return ""
	}
}

func (result *Result) PreviewURL(res string) *url.URL {
	fileParts := strings.SplitN(result.Image, ".", 2)
	cfg := resources.Resources[res]

	previewURL := url.URL{
		Scheme: cfg.GetString("scheme"),
		Host:   cfg.GetString("host"),
	}

	var fileName string
	switch cfg.GetString("thumbs.name") {
	case "images":
		fileName = fileParts[0]
	default:
		fileName = result.Hash
	}

	var fileFormat string
	switch cfg.GetString("thumbs.format") {
	case "images":
		fileFormat = fileParts[1]
	default:
		fileFormat = cfg.GetString("thumbs.format")
		if fileFormat == "" {
			fileFormat = fileParts[1]
		}
	}

	previewURL.Path = fmt.Sprint(
		cfg.GetString("thumbs.dir"),
		result.Directory,
		cfg.GetString("thumbs.part"),
		fileName,
		".",
		fileFormat,
	)

	return &previewURL
}

func (result *Result) SampleURL(resource string) *url.URL {
	cfg := resources.Resources[resource]
	sampleURL := url.URL{
		Scheme: cfg.GetString("scheme"),
		Host:   cfg.GetString("host"),
	}

	fileParts := strings.SplitN(result.Image, ".", 2)

	var fileName string
	switch cfg.GetString("samples.name") {
	case "thumbs":
		fileName = result.Hash
	default:
		fileName = fileParts[0]
	}

	var fileFormat string
	switch cfg.GetString("samples.format") {
	case "images":
		fileFormat = fileParts[1]
	default:
		fileFormat = cfg.GetString("samples.format")
		if fileFormat == "" {
			fileFormat = fileParts[1]
		}
	}

	sampleURL.Path = fmt.Sprint(
		cfg.GetString("samples.dir"),
		result.Directory,
		cfg.GetString("samples.part"),
		fileName,
		".",
		fileFormat,
	)

	return &sampleURL
}

func (result *Result) FileURL(resource string) *url.URL {
	cfg := resources.Resources[resource]
	fileURL := url.URL{
		Scheme: cfg.GetString("scheme"),
		Host:   cfg.GetString("host"),
	}

	fileParts := strings.SplitN(result.Image, ".", 2)
	fileFormat := fileParts[1]

	var fileName string
	switch cfg.GetString("images.name") {
	case "thumbs":
		fileName = result.Hash
	default:
		fileName = fileParts[0]
	}

	fileURL.Path = fmt.Sprint(
		cfg.GetString("images.dir"),
		result.Directory,
		cfg.GetString("images.part"),
		fileName,
		".",
		fileFormat,
	)

	return &fileURL
}
