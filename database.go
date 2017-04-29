package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	bolt "github.com/boltdb/bolt"
)

type (
	User struct {
		ID                          int
		Language                    string
		Roles, Blacklist, Whitelist []string
		Ratings                     Ratings
		Patreon                     Patreon
	}

	Ratings struct {
		Safe, Questionale, Explicit bool
	}

	Patreon struct {
		FullName, AccessToken, RefreshToken string
	}
)

var (
	bktPatreon   = []byte("patreon")
	bktRatings   = []byte("ratings")
	bktRoles     = []byte("roles")
	bktBlacklist = []byte("black")
	bktWhitelist = []byte("white")

	db *bolt.DB
)

func dbInit() {
	db, err = bolt.Open(*flagDB, 0600, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func createUser(userID int) (*User, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		user, err := tx.CreateBucket([]byte(strconv.Itoa(userID)))
		if err != nil {
			return err
		}

		user.Put([]byte("id"), []byte(strconv.Itoa(userID)))
		user.Put([]byte("language"), []byte("en"))

		black, err := user.CreateBucket(bktBlacklist)
		if err != nil {
			return err
		}

		black.Put([]byte("loli*"), nil)
		black.Put([]byte("scat"), nil)

		if _, err := user.CreateBucket(bktWhitelist); err != nil {
			return err
		}

		roles, err := user.CreateBucket(bktRoles)
		if err != nil {
			return err
		}

		if userID == adm {
			roles.Put([]byte("admin"), nil)
			roles.Put([]byte("patron"), nil)
		}

		roles.Put([]byte("user"), nil)

		ratings, err := user.CreateBucket(bktRatings)
		if err != nil {
			return err
		}

		ratings.Put([]byte("safe"), strconv.AppendBool(nil, true))
		ratings.Put([]byte("questionable"), strconv.AppendBool(nil, false))
		ratings.Put([]byte("explicit"), strconv.AppendBool(nil, false))

		_, err = user.CreateBucket(bktPatreon)
		return err
	}); err != nil {
		return nil, err
	}
	return getUser(userID)
}

func (u *User) changeLanguage(lang string) (*User, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(u.ID)))
		if user == nil {
			return errors.New("user not exist")
		}

		return user.Put([]byte("language"), []byte(lang))
	}); err != nil {
		return nil, err
	}

	return getUser(u.ID)
}

func (u *User) addRoles(roles ...string) (*User, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(u.ID)))
		if user == nil {
			return errors.New("user not exist")
		}

		userRoles := user.Bucket(bktRoles)

		for _, role := range roles {
			if err := userRoles.Put([]byte(role), nil); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return getUser(u.ID)
}

func (u *User) removeRoles(roles ...string) (*User, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(u.ID)))
		if user == nil {
			return errors.New("user not exist")
		}

		userRoles := user.Bucket(bktRoles)

		for _, role := range roles {
			if err := userRoles.Delete([]byte(role)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return getUser(u.ID)
}

func (u *User) changeRatings(safe, questionable, explicit bool) (*User, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(u.ID)))
		if user == nil {
			return errors.New("user not exist")
		}

		ratings := user.Bucket(bktRatings)
		if err := ratings.Put([]byte("safe"), strconv.AppendBool(nil, safe)); err != nil {
			return err
		}

		if err := ratings.Put([]byte("questionable"), strconv.AppendBool(nil, questionable)); err != nil {
			return err
		}

		return ratings.Put([]byte("explicit"), strconv.AppendBool(nil, explicit))
	}); err != nil {
		return nil, err
	}
	return getUser(u.ID)
}

func (u *User) patreonSave(name, access, refresh string) (*User, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(u.ID)))
		if user == nil {
			return errors.New("user not exist")
		}
		pat := user.Bucket(bktPatreon)

		if err := pat.Put([]byte("name"), []byte(name)); err != nil {
			return err
		}

		if err := pat.Put([]byte("access"), []byte(access)); err != nil {
			return err
		}

		return pat.Put([]byte("refresh"), []byte(refresh))
	}); err != nil {
		return nil, err
	}

	return getUser(u.ID)
}

func (u *User) tagsRewrite(black bool, tags []string) error {
	return db.Update(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(u.ID)))
		if user == nil {
			return errors.New("user not exist")
		}

		var bkt *bolt.Bucket
		if black {
			err := user.DeleteBucket(bktBlacklist)
			if err != nil {
				return err
			}

			bkt, err = user.CreateBucket(bktBlacklist)
			if err != nil {
				return err
			}
		} else {
			err := user.DeleteBucket(bktWhitelist)
			if err != nil {
				return err
			}

			bkt, err = user.CreateBucket(bktWhitelist)
			if err != nil {
				return err
			}
		}

		for _, tag := range tags {
			tag = strings.ToLower(strings.TrimLeft(tag, "-"))
			bkt.Put([]byte(tag), nil)
		}
		return nil
	})
}

func (u *User) tagRemove(black bool, tag string) (*User, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(u.ID)))
		if user == nil {
			return errors.New("user not exist")
		}

		var bkt *bolt.Bucket
		if black {
			bkt = user.Bucket(bktBlacklist)
		} else {
			bkt = user.Bucket(bktWhitelist)
		}

		return bkt.Delete([]byte(tag))
	}); err != nil {
		return createUser(u.ID)
	}

	return getUser(u.ID)
}

func getUser(userID int) (*User, error) {
	var usr User
	if err := db.View(func(tx *bolt.Tx) error {
		user := tx.Bucket([]byte(strconv.Itoa(userID)))
		if user == nil {
			return errors.New("user not exist")
		}

		usr.Language = string(user.Get([]byte("language")))
		usr.ID, err = strconv.Atoi(string(user.Get([]byte("id"))))
		if err != nil {
			return err
		}

		black := user.Bucket(bktBlacklist)
		var list []string
		if err := black.ForEach(func(key, val []byte) error {
			list = append(list, string(key))
			return nil
		}); err != nil {
			return err
		}
		usr.Blacklist = list

		white := user.Bucket(bktWhitelist)
		list = nil
		if err := white.ForEach(func(key, val []byte) error {
			list = append(list, string(key))
			return nil
		}); err != nil {
			return err
		}
		usr.Whitelist = list

		roles := user.Bucket(bktRoles)
		list = nil
		if err := roles.ForEach(func(key, val []byte) error {
			list = append(list, string(key))
			return nil
		}); err != nil {
			return err
		}
		usr.Roles = list

		ratings := user.Bucket(bktRatings)
		var uRatings Ratings
		uRatings.Safe, err = strconv.ParseBool(string(ratings.Get([]byte("safe"))))
		if err != nil {
			return err
		}
		uRatings.Questionale, err = strconv.ParseBool(string(ratings.Get([]byte("questionable"))))
		if err != nil {
			return err
		}
		uRatings.Explicit, err = strconv.ParseBool(string(ratings.Get([]byte("explicit"))))
		if err != nil {
			return err
		}
		usr.Ratings = uRatings

		patreon := user.Bucket(bktPatreon)
		if patreon != nil {
			var uPatreon Patreon
			uPatreon.FullName = string(patreon.Get([]byte("name")))
			uPatreon.AccessToken = string(patreon.Get([]byte("access")))
			uPatreon.RefreshToken = string(patreon.Get([]byte("refresh")))
			usr.Patreon = uPatreon
		}

		return nil
	}); err != nil {
		return createUser(userID)
	}
	return &usr, nil
}

func getUsers() ([]User, error) {
	var users []User
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, bkt *bolt.Bucket) error {
			id, err := strconv.Atoi(string(name))
			if err != nil {
				return err
			}

			usr, err := getUser(id)
			if err != nil {
				return err
			}

			users = append(users, *usr)
			return nil
		})
	})
	return users, err
}
