package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

func dbGetUser(id int) (*user, error) {
	log.Ln("db:get:user")

	var usr user
	usr.ID = id
	usr.Resources = make(map[string]bool)
	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		// Language
		usr.Language, err = tx.Get(fmt.Sprint("user:", id, ":language"))
		if err != nil {
			return err
		}

		// Roles
		roleUser, _ := tx.Get(fmt.Sprint("user:", id, ":role:user"))
		usr.Roles.User, _ = strconv.ParseBool(roleUser)

		rolePatron, _ := tx.Get(fmt.Sprint("user:", id, ":role:patron"))
		usr.Roles.Patron, _ = strconv.ParseBool(rolePatron)

		roleManager, _ := tx.Get(fmt.Sprint("user:", id, ":role:manager"))
		usr.Roles.Manager, _ = strconv.ParseBool(roleManager)

		roleAdmin, _ := tx.Get(fmt.Sprint("user:", id, ":role:admin"))
		usr.Roles.Admin, _ = strconv.ParseBool(roleAdmin)

		// Ratings
		rateSafe, _ := tx.Get(fmt.Sprint("user:", id, ":rating:safe"))
		usr.Ratings.Safe, _ = strconv.ParseBool(rateSafe)

		rateQuestionable, _ := tx.Get(
			fmt.Sprint("user:", id, ":rating:questionable"),
		)
		usr.Ratings.Questionable, _ = strconv.ParseBool(rateQuestionable)

		rateExplicit, _ := tx.Get(fmt.Sprint("user:", id, ":rating:explicit"))
		usr.Ratings.Exlplicit, _ = strconv.ParseBool(rateExplicit)

		// Resources
		tx.AscendKeys(
			fmt.Sprint("user:", id, ":resource:*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				name := keys[3]
				usr.Resources[name], _ = strconv.ParseBool(val)
				return true
			},
		)

		// Blacklist
		tx.AscendKeys(
			fmt.Sprint("user:", id, ":", blackList, ":*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				usr.Blacklist = append(usr.Blacklist, keys[3])
				return true
			},
		)

		// Whitelist
		tx.AscendKeys(
			fmt.Sprint("user:", id, ":", whiteList, ":*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				usr.Whitelist = append(usr.Whitelist, keys[3])
				return true
			},
		)

		return nil
	})

	sort.Strings(usr.Whitelist)
	sort.Strings(usr.Blacklist)

	return &usr, err
}
