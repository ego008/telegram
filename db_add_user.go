package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

func dbAddUser(id int, lang string) (*user, error) {
	log.Ln("db:add:user")

	var usr user
	usr.Blacklist = []string{
		"loli*",
		"scat",
		"uncensored",
	}
	usr.ID = id
	usr.Language = lang
	usr.Ratings = ratings{
		Safe:         true,
		Questionable: false,
		Exlplicit:    false,
	}
	usr.Resources = make(map[string]bool)
	usr.Roles = roles{
		User:    false,
		Patron:  false,
		Manager: false,
		Admin:   false,
	}
	usr.Whitelist = make([]string, 0)

	for k, _ := range resources {
		usr.Resources[k] = false
	}

	err := db.Update(func(tx *buntdb.Tx) error {
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
