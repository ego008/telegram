package db

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/HentaiDB/HentaiDBot/internal/models"
	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

func User(id int) (*models.User, error) {
	log.Ln("db:get:user")

	usr := models.User{
		ID:        id,
		Resources: make(map[string]bool),
	}

	err := DB.View(func(tx *buntdb.Tx) error {
		var err error
		// Language
		usr.Language, err = tx.Get(fmt.Sprint("user:", id, ":language"))
		if err != nil {
			return err
		}

		// Roles
		roleUser, err := tx.Get(fmt.Sprint("user:", id, ":role:user"))
		if err != nil {
			return err
		}
		usr.Roles.User, err = strconv.ParseBool(roleUser)
		if err != nil {
			return err
		}

		rolePatron, err := tx.Get(fmt.Sprint("user:", id, ":role:patron"))
		if err != nil {
			return err
		}
		usr.Roles.Patron, err = strconv.ParseBool(rolePatron)
		if err != nil {
			return err
		}

		roleManager, err := tx.Get(fmt.Sprint("user:", id, ":role:manager"))
		if err != nil {
			return err
		}
		usr.Roles.Manager, err = strconv.ParseBool(roleManager)
		if err != nil {
			return err
		}

		roleAdmin, err := tx.Get(fmt.Sprint("user:", id, ":role:admin"))
		if err != nil {
			return err
		}
		usr.Roles.Admin, err = strconv.ParseBool(roleAdmin)
		if err != nil {
			return err
		}

		// Ratings
		rateSafe, err := tx.Get(fmt.Sprint("user:", id, ":rating:safe"))
		if err != nil {
			return err
		}
		usr.Ratings.Safe, err = strconv.ParseBool(rateSafe)
		if err != nil {
			return err
		}

		rateQuestionable, err := tx.Get(
			fmt.Sprint("user:", id, ":rating:questionable"),
		)
		if err != nil {
			return err
		}
		usr.Ratings.Questionable, err = strconv.ParseBool(rateQuestionable)
		if err != nil {
			return err
		}

		rateExplicit, err := tx.Get(fmt.Sprint("user:", id, ":rating:explicit"))
		if err != nil {
			return err
		}
		usr.Ratings.Exlplicit, err = strconv.ParseBool(rateExplicit)
		if err != nil {
			return err
		}

		// Types
		typeImage, err := tx.Get(fmt.Sprint("user:", id, ":type:image"))
		if err != nil {
			return err
		}
		usr.ContentTypes.Image, err = strconv.ParseBool(typeImage)
		if err != nil {
			return err
		}

		typeAnimation, err := tx.Get(fmt.Sprint("user:", id, ":type:animation"))
		if err != nil {
			return err
		}
		usr.ContentTypes.Animation, err = strconv.ParseBool(typeAnimation)
		if err != nil {
			return err
		}

		typeVideo, err := tx.Get(fmt.Sprint("user:", id, ":type:video"))
		if err != nil {
			return err
		}
		usr.ContentTypes.Video, err = strconv.ParseBool(typeVideo)
		if err != nil {
			return err
		}

		// Resources
		if err = tx.AscendKeys(
			fmt.Sprint("user:", id, ":resource:*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				name := keys[3]
				usr.Resources[name], err = strconv.ParseBool(val)
				if err != nil {
					return true
				}
				return true
			},
		); err != nil {
			return err
		}

		// Blacklist
		if err = tx.AscendKeys(
			fmt.Sprint("user:", id, ":", blackList, ":*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				usr.Blacklist = append(usr.Blacklist, keys[3])
				return true
			},
		); err != nil {
			return err
		}

		// Whitelist
		err = tx.AscendKeys(
			fmt.Sprint("user:", id, ":", whiteList, ":*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				usr.Whitelist = append(usr.Whitelist, keys[3])
				return true
			},
		)

		return err
	})
	if err != nil {
		return &usr, err
	}

	sort.Strings(usr.Whitelist)
	sort.Strings(usr.Blacklist)

	return &usr, nil
}
