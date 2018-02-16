package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/models"
	"github.com/HentaiDB/HentaiDBot/internal/resources"
	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

func AddUser(id int, lang string) (*models.User, error) {
	log.Ln("db:add:user")

	usr := models.User{
		Blacklist: []string{
			"loli*",
			"scat",
			"uncensored",
		},
		ID:       id,
		Language: lang,
		Ratings: models.Ratings{
			Safe:         true,
			Questionable: false,
			Exlplicit:    false,
		},
		Resources: make(map[string]bool),
		Roles: models.Roles{
			User:    false,
			Patron:  false,
			Manager: false,
			Admin:   false,
		},
		ContentTypes: models.ContentTypes{
			Image:     true,
			Animation: true,
			Video:     true,
		},
		Whitelist: make([]string, 0),
	}

	for k, _ := range resources.Resources {
		usr.Resources[k] = false
	}

	err := DB.Update(func(tx *buntdb.Tx) error {
		// Language
		tx.Set(
			fmt.Sprint("user:", id, ":language"),
			usr.Language, nil,
		)

		// Roles
		tx.Set(
			fmt.Sprint("user:", id, ":role:user"),
			strconv.FormatBool(usr.Roles.User), nil,
		)
		tx.Set(
			fmt.Sprint("user:", id, ":role:patron"),
			strconv.FormatBool(usr.Roles.Patron), nil,
		)
		tx.Set(
			fmt.Sprint("user:", id, ":role:manager"),
			strconv.FormatBool(usr.Roles.Manager), nil,
		)
		tx.Set(
			fmt.Sprint("user:", id, ":role:admin"),
			strconv.FormatBool(usr.Roles.Admin), nil,
		)

		// Ratings
		tx.Set(
			fmt.Sprint("user:", id, ":rating:safe"),
			strconv.FormatBool(usr.Ratings.Safe), nil,
		)
		tx.Set(
			fmt.Sprint("user:", id, ":rating:questionable"),
			strconv.FormatBool(usr.Ratings.Questionable), nil,
		)
		tx.Set(
			fmt.Sprint("user:", id, ":rating:explicit"),
			strconv.FormatBool(usr.Ratings.Exlplicit), nil,
		)

		// Types
		tx.Set(
			fmt.Sprint("user:", id, ":type:image"),
			strconv.FormatBool(usr.ContentTypes.Image), nil,
		)
		tx.Set(
			fmt.Sprint("user:", id, ":type:animation"),
			strconv.FormatBool(usr.ContentTypes.Animation), nil,
		)
		tx.Set(
			fmt.Sprint("user:", id, ":type:video"),
			strconv.FormatBool(usr.ContentTypes.Video), nil,
		)

		// Resources
		for k, v := range usr.Resources {
			tx.Set(
				fmt.Sprint("user:", id, ":resource:", k),
				strconv.FormatBool(v), nil,
			)
		}

		// Blacklist
		for i := range usr.Blacklist {
			tag := strings.ToLower(usr.Blacklist[i])
			tx.Set(
				fmt.Sprint("user:", id, ":", blackList, ":", tag),
				"", nil,
			)
		}

		// Whitelist
		for i := range usr.Whitelist {
			tag := strings.ToLower(usr.Whitelist[i])
			tx.Set(
				fmt.Sprint("user:", id, ":", whiteList, ":", tag),
				"", nil,
			)
		}

		return nil
	})

	return &usr, err
}
