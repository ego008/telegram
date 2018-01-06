package main

import (
	"fmt"
	"sort"
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

func (usr *user) toggleRatingSafe() error {
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

func (usr *user) toggleRatingQuestionable() error {
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

func (usr *user) toggleRatingExplicit() error {
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

func (usr *user) toggleTypeImage() error {
	log.Ln("db:user:toggle:image")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":type:image"),
			strconv.FormatBool(!usr.Types.Image), nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Types.Image = !usr.Types.Image
	return nil
}

func (usr *user) toggleTypeAnimation() error {
	log.Ln("db:user:toggle:animation")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":type:animation"),
			strconv.FormatBool(!usr.Types.Animation), nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Types.Animation = !usr.Types.Animation
	return nil
}

func (usr *user) toggleTypeVideo() error {
	log.Ln("db:user:toggle:video")

	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":type:video"),
			strconv.FormatBool(!usr.Types.Video), nil,
		)
		return err
	})
	if err != nil {
		return err
	}

	usr.Types.Video = !usr.Types.Video
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

func (usr *user) addListTags(listType string, tags ...string) error {
	log.Ln("db:user:add:" + listType + ":tags")

	err := db.Update(func(tx *buntdb.Tx) error {
		for i := range tags {
			tag := strings.ToLower(tags[i])
			_, _, err := tx.Set(
				fmt.Sprint("user:", usr.ID, ":", listType, ":", tag), "", nil,
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

	switch listType {
	case whiteList:
		for i := range tags {
			for j := range usr.Whitelist {
				if usr.Whitelist[j] == tags[i] {
					continue
				}

				usr.Whitelist = append(usr.Whitelist, tags[i])
			}
		}
	case blackList:
		for i := range tags {
			for j := range usr.Blacklist {
				if usr.Blacklist[j] == tags[i] {
					continue
				}

				usr.Blacklist = append(usr.Blacklist, tags[i])
			}
		}
	}

	return nil
}

func (usr *user) removeListTag(listType string, tag string) error {
	log.Ln("db:user:remove:" + listType + ":tags")
	pattern := fmt.Sprint("user:", usr.ID, ":", listType, ":")

	err := db.Update(func(tx *buntdb.Tx) error {
		var tagKey string
		err := tx.AscendKeys(
			fmt.Sprint(pattern, "*"), func(key, val string) bool {
				if strings.TrimPrefix(key, pattern) == tag {
					tagKey = key
					return false
				}
				return true
			},
		)
		if err != nil {
			return err
		}

		_, err = tx.Delete(tagKey)
		return err
	})
	if err != nil {
		return err
	}

	switch listType {
	case whiteList:
		for i := range usr.Whitelist {
			if usr.Whitelist[i] != tag {
				continue
			}

			usr.Whitelist = append(usr.Whitelist[:i], usr.Whitelist[i+1:]...)
			break
		}
		sort.Strings(usr.Whitelist)
	case blackList:
		for i := range usr.Blacklist {
			if usr.Blacklist[i] != tag {
				continue
			}

			usr.Blacklist = append(usr.Blacklist[:i], usr.Blacklist[i+1:]...)
			break
		}
		sort.Strings(usr.Blacklist)
	}

	return nil
}
