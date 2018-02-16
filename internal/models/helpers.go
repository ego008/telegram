package models

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/i18n"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
)

func (usr *User) GetRatingsStatus() (string, error) {
	T, err := i18n.SwitchTo(usr.Language)
	if err != nil {
		return "", err
	}

	switch {
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe only
		return T("rating_safe"), nil
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// questionable only
		return T("rating_questionable"), nil
	case !usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// explicit only
		return T("rating_explicit"), nil
	case usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe + questionable
		return fmt.Sprint(T("rating_safe"), " + ", T("rating_questionable")), nil
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// safe + explicit
		return fmt.Sprint(T("rating_safe"), " + ", T("rating_explicit")), nil
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// questionable + explicit
		return fmt.Sprint(T("rating_questionable"), " + ", T("rating_explicit")), nil
	default:
		// all ratings enabled/diabled
		return T("rating_all"), nil
	}
}

func (usr *User) GetRatingsFilter() string {
	switch {
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe only
		return "rating:safe"
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// questionable only
		return "rating:questionable"
	case !usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// explicit only
		return "rating:explicit"
	case usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe + questionable
		return "-rating:explicit"
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// safe + explicit
		return "-rating:questionable"
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// questionable + explicit
		return "-rating:safe"
	default:
		// all ratings enabled/diabled
		return ""
	}
}

func CheckInterface(src interface{}) string {
	if src != nil {
		return src.(string)
	}
	return ""
}

func (result *Result) PreviewURL(resource string) *url.URL {
	var previewURL url.URL
	res := resources.Resources[resource]

	previewURL.Scheme = res.UString("scheme", "http")
	previewURL.Host = res.UString("host")

	fileParts := strings.SplitN(result.Image, ".", 2)

	var fileName string
	switch res.UString("thumbs.name") {
	case "images":
		fileName = fileParts[0]
	default:
		fileName = result.Hash
	}

	var fileFormat string
	switch res.UString("thumbs.format") {
	case "images":
		fileFormat = fileParts[1]
	default:
		fileFormat = res.UString("thumbs.format", fileParts[1])
	}

	previewURL.Path = fmt.Sprint(
		res.UString("thumbs.dir"),
		result.Directory,
		res.UString("thumbs.part"),
		fileName,
		".",
		fileFormat,
	)

	return &previewURL
}

func (result *Result) SampleURL(resource string) *url.URL {
	var sampleURL url.URL
	res := resources.Resources[resource]

	sampleURL.Scheme = res.UString("scheme", "http")
	sampleURL.Host = res.UString("host")

	fileParts := strings.SplitN(result.Image, ".", 2)

	var fileName string
	switch res.UString("samples.name") {
	case "thumbs":
		fileName = result.Hash
	default:
		fileName = fileParts[0]
	}

	var fileFormat string
	switch res.UString("samples.format") {
	case "images":
		fileFormat = fileParts[1]
	default:
		fileFormat = res.UString("samples.format", fileParts[1])
	}

	sampleURL.Path = fmt.Sprint(
		res.UString("samples.dir"),
		result.Directory,
		res.UString("samples.part"),
		fileName,
		".",
		fileFormat,
	)

	return &sampleURL
}

func (result *Result) FileURL(resource string) *url.URL {
	var fileURL url.URL
	res := resources.Resources[resource]

	fileURL.Scheme = res.UString("scheme", "http")
	fileURL.Host = res.UString("host")

	fileParts := strings.SplitN(result.Image, ".", 2)

	fileFormat := fileParts[1]

	var fileName string
	switch res.UString("images.name") {
	case "thumbs":
		fileName = result.Hash
	default:
		fileName = fileParts[0]
	}

	fileURL.Path = fmt.Sprint(
		res.UString("images.dir"),
		result.Directory,
		res.UString("images.part"),
		fileName,
		".",
		fileFormat,
	)

	return &fileURL
}
