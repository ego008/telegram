package main

import (
	"fmt"
	"net/url"
	"strings"
)

func (usr *user) getRatingsStatus() (string, error) {
	T, err := langSwitch(usr.Language)
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

func (usr *user) getRatingsFilter() string {
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

func checkInterface(src interface{}) string {
	if src != nil {
		return src.(string)
	}
	return ""
}

func (post *post) previewURL(resName string) *url.URL {
	var previewURL url.URL
	resource := resources[resName]

	previewURL.Scheme = resource.UString("scheme", "http")
	previewURL.Host = resource.UString("host")

	fileParts := strings.SplitN(post.Image, ".", 2)

	var fileName string
	switch resource.UString("thumbs.name") {
	case "images":
		fileName = fileParts[0]
	default:
		fileName = post.Hash
	}

	var fileFormat string
	switch resource.UString("thumbs.format") {
	case "images":
		fileFormat = fileParts[1]
	default:
		fileFormat = resource.UString("thumbs.format", fileParts[1])
	}

	previewURL.Path = fmt.Sprint(
		resource.UString("thumbs.dir"),
		post.Directory,
		resource.UString("thumbs.part"),
		fileName,
		".",
		fileFormat,
	)

	return &previewURL
}

func (post *post) sampleURL(resName string) *url.URL {
	var sampleURL url.URL
	resource := resources[resName]

	sampleURL.Scheme = resource.UString("scheme", "http")
	sampleURL.Host = resource.UString("host")

	fileParts := strings.SplitN(post.Image, ".", 2)

	var fileName string
	switch resource.UString("samples.name") {
	case "thumbs":
		fileName = post.Hash
	default:
		fileName = fileParts[0]
	}

	var fileFormat string
	switch resource.UString("samples.format") {
	case "images":
		fileFormat = fileParts[1]
	default:
		fileFormat = resource.UString("samples.format", fileParts[1])
	}

	sampleURL.Path = fmt.Sprint(
		resource.UString("samples.dir"),
		post.Directory,
		resource.UString("samples.part"),
		fileName,
		".",
		fileFormat,
	)

	return &sampleURL
}

func (post *post) fileURL(resName string) *url.URL {
	var fileURL url.URL
	resource := resources[resName]

	fileURL.Scheme = resource.UString("scheme", "http")
	fileURL.Host = resource.UString("host")

	fileParts := strings.SplitN(post.Image, ".", 2)

	fileFormat := fileParts[1]

	var fileName string
	switch resource.UString("images.name") {
	case "thumbs":
		fileName = post.Hash
	default:
		fileName = fileParts[0]
	}

	fileURL.Path = fmt.Sprint(
		resource.UString("images.dir"),
		post.Directory,
		resource.UString("images.part"),
		fileName,
		".",
		fileFormat,
	)

	return &fileURL
}
