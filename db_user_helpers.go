package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

func dbGetUserElseAdd(id int, lang string) (usr *user, err error) {
	usr, err = dbGetUser(id)
	if err == errNotFound {
		usr, err = dbAddUser(id, lang)
	}

	return
}

func (usr *user) Destroy() error {
	log.Ln("db:user:destroy")

	return db.Update(func(tx *buntdb.Tx) error {
		return tx.Ascend("users", func(key, val string) bool {
			keys := strings.Split(key, ":")
			if keys[1] != strconv.Itoa(usr.ID) {
				return true
			}

			tx.Delete(key)
			return true
		})
	})
}

func (usr *user) toggleSafe() error {
	log.Ln("db:user:toggle:safe")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":rating:safe"),
			strconv.FormatBool(!usr.Ratings.Safe), nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Ratings.Safe = !usr.Ratings.Safe
	return nil
}

func (usr *user) toggleQuestionable() error {
	log.Ln("db:user:toggle:questionable")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":rating:questionable"),
			strconv.FormatBool(!usr.Ratings.Questionable), nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Ratings.Questionable = !usr.Ratings.Questionable
	return nil
}

func (usr *user) toggleExplicit() error {
	log.Ln("db:user:toggle:explicit")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":rating:explicit"),
			strconv.FormatBool(!usr.Ratings.Exlplicit), nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Ratings.Exlplicit = !usr.Ratings.Exlplicit
	return nil
}

func (usr *user) toggleResource(res string) error {
	log.Ln("db:user:toggle:resource")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":resource:", res),
			strconv.FormatBool(!usr.Resources[res]), nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Resources[res] = !usr.Resources[res]
	return nil
}

func (usr *user) setLanguage(lang string) error {
	log.Ln("db:user:set:language")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":language"),
			lang, nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Language = lang
	return nil
}

func (usr *user) addListTags(listType string, tag ...string) error {
	log.Ln("db:user:add:" + listType + ":tags")

	err := db.Update(func(tx *buntdb.Tx) error {
		for i := range tag {
			tag := strings.ToLower(tag[i])
			_, _, err := tx.Set(
				fmt.Sprint("user:", usr.ID, ":", listType, ":", tag),
				"", nil,
			)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	var tags []string
	err = db.View(func(tx *buntdb.Tx) error {
		pattern := fmt.Sprint("user:", usr.ID, ":", listType, ":")
		err := tx.AscendKeys(
			fmt.Sprint(pattern, "*"),
			func(key, val string) bool {
				tags = append(tags, strings.TrimPrefix(key, pattern))
				return true
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	switch listType {
	case whiteList:
		usr.Whitelist = tags
	case blackList:
		usr.Blacklist = tags
	}

	return nil
}

func (usr *user) removeListTag(listType string, tag string) error {
	log.Ln("db:user:remove:" + listType + ":tags")
	pattern := fmt.Sprint("user:", usr.ID, ":", listType, ":")

	err := db.Update(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint(pattern, "*"), func(key, val string) bool {
				if strings.TrimPrefix(key, pattern) == tag {
					tx.Delete(key)
					return false
				}
				return true
			},
		)
	})
	if err != nil {
		return err
	}

	var tags []string
	err = db.View(func(tx *buntdb.Tx) error {
		pattern := fmt.Sprint("user:", usr.ID, ":", listType, ":")
		err := tx.AscendKeys(
			fmt.Sprint(pattern, "*"),
			func(key, val string) bool {
				tags = append(tags, strings.TrimPrefix(key, pattern))
				return true
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	switch listType {
	case whiteList:
		usr.Whitelist = tags
	case blackList:
		usr.Blacklist = tags
	}

	return nil
}
