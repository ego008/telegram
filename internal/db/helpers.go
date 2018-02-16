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

func Destroy(userID int) error {
	log.Ln("db:user:destroy")

	return DB.Update(func(tx *buntdb.Tx) error {
		return tx.Ascend("users", func(key, val string) bool {
			keys := strings.Split(key, ":")
			if keys[1] != strconv.Itoa(userID) {
				return true
			}

			tx.Delete(key)
			return true
		})
	})
}

func ToggleRatingSafe(usr *models.User) error {
	log.Ln("db:user:toggle:safe")

	if err := DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":rating:safe"),
			strconv.FormatBool(!usr.Ratings.Safe), nil,
		)
		return err
	}); err != nil {
		return err
	}

	usr.Ratings.Safe = !usr.Ratings.Safe
	return nil
}

func ToggleRatingQuestionable(usr *models.User) error {
	log.Ln("db:user:toggle:questionable")

	err := DB.Update(func(tx *buntdb.Tx) error {
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

func ToggleRatingExplicit(usr *models.User) error {
	log.Ln("db:user:toggle:explicit")

	if err := DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":rating:explicit"),
			strconv.FormatBool(!usr.Ratings.Exlplicit), nil,
		)
		return err
	}); err != nil {
		return err
	}

	usr.Ratings.Exlplicit = !usr.Ratings.Exlplicit
	return nil
}

func ToggleTypeImage(usr *models.User) error {
	log.Ln("db:user:toggle:image")

	if err := DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":type:image"),
			strconv.FormatBool(!usr.ContentTypes.Image), nil,
		)
		return err
	}); err != nil {
		return err
	}

	usr.ContentTypes.Image = !usr.ContentTypes.Image
	return nil
}

func ToggleTypeAnimation(usr *models.User) error {
	log.Ln("db:user:toggle:animation")

	if err := DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":type:animation"),
			strconv.FormatBool(!usr.ContentTypes.Animation), nil,
		)
		return err
	}); err != nil {
		return err
	}

	usr.ContentTypes.Animation = !usr.ContentTypes.Animation
	return nil
}

func ToggleTypeVideo(usr *models.User) error {
	log.Ln("db:user:toggle:video")

	if err := DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":type:video"),
			strconv.FormatBool(!usr.ContentTypes.Video), nil,
		)
		return err
	}); err != nil {
		return err
	}

	usr.ContentTypes.Video = !usr.ContentTypes.Video
	return nil
}

func ToggleResource(usr *models.User, resource string) error {
	log.Ln("db:user:toggle:resource")

	if err := DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":resource:", resource),
			strconv.FormatBool(!usr.Resources[resource]), nil,
		)
		return err
	}); err != nil {
		return err
	}

	usr.Resources[resource] = !usr.Resources[resource]
	return nil
}

func SetLanguage(usr *models.User, lang string) error {
	log.Ln("db:user:set:language")

	if err := DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", usr.ID, ":language"),
			lang, nil,
		)
		return err
	}); err != nil {
		return err
	}

	usr.Language = lang
	return nil
}

func AddListTags(usr *models.User, listType string, tags ...string) error {
	log.Ln("db:user:add:" + listType + ":tags")

	if err := DB.Update(func(tx *buntdb.Tx) error {
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
	}); err != nil {
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

func RemoveListTag(usr *models.User, listType string, tag string) error {
	log.Ln("db:user:remove:" + listType + ":tags")
	pattern := fmt.Sprint("user:", usr.ID, ":", listType, ":")

	if err := DB.Update(func(tx *buntdb.Tx) error {
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
	}); err != nil {
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
