package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	// tg "gopkg.in/telegram-bot-api.v4"
	// log "github.com/kirillDanshin/dlog"
)

type UserDB struct {
	Language string
	NSFW     bool
	// Menu     string
	Role string
	Hits int
}

var (
	db       *bolt.DB
	bktUsers = []byte("users")
)

func init() {
	go func() {
		var err error
		db, err = bolt.Open("hentai.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bktUsers)
			return err
		}); err != nil {
			panic(err.Error())
		}

		select {}
	}()
}

func CreateUserBD(id int) error {
	return db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.Bucket(bktUsers).CreateBucket([]byte(strconv.Itoa(id)))
		if err != nil {
			return err
		}

		for _, admin := range cfg["admins"].([]interface{}) {
			if id == int(admin.(float64)) {
				bkt.Put([]byte("role"), []byte("anon"))
			} else {
				for _, patron := range cfg["patrons"].([]interface{}) {
					if id == int(patron.(float64)) {
						bkt.Put([]byte("role"), []byte("patron"))
					} else {
						bkt.Put([]byte("role"), []byte("anon"))
					}
				}
			}
		}

		bkt.Put([]byte("lang"), []byte("en-us"))
		bkt.Put([]byte("nsfw"), strconv.AppendBool(nil, false))
		// bkt.Put([]byte("menu"), []byte("start"))
		bkt.Put([]byte("hits"), []byte(strconv.Itoa(0)))
		return nil
	})
}

func ChangeLangBD(id int, lang string) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		return bkt.Put([]byte("lang"), []byte(lang))
	}); err != nil {
		CreateUserBD(id)
		return ChangeLangBD(id, lang)
	}
	return nil
}

func ChangeRoleBD(id int, role string) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		return bkt.Put([]byte("role"), []byte(role))
	}); err != nil {
		CreateUserBD(id)
		return ChangeRoleBD(id, role)
	}
	return nil
}

func ChangeFilterDB(id int, nsfw bool) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		return bkt.Put([]byte("nsfw"), strconv.AppendBool(nil, nsfw))
	}); err != nil {
		CreateUserBD(id)
		return ChangeFilterDB(id, nsfw)
	}
	return nil
}

func GetUserDB(id int) (*UserDB, error) {
	var usr UserDB
	if err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		usr.Language = string(bkt.Get([]byte("lang")))
		usr.NSFW, _ = strconv.ParseBool(string(bkt.Get([]byte("nsfw"))))
		// usr.Menu = string(bkt.Get([]byte("menu")))
		usr.Role = string(bkt.Get([]byte("role")))
		usr.Hits, _ = strconv.Atoi(string(bkt.Get([]byte("hits"))))
		return nil
	}); err != nil {
		CreateUserBD(id)
		return GetUserDB(id)
	}
	return &usr, nil
}

func AddHitsDB(id int) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		hits, _ := strconv.Atoi(string(bkt.Get([]byte("hits"))))
		hits++
		return bkt.Put([]byte("hits"), []byte(strconv.Itoa(hits)))
	}); err != nil {
		CreateUserBD(id)
		return AddHitsDB(id)
	}
	return nil
}
