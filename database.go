package main

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	// tg "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/kirillDanshin/dlog"
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
	bktPosts = []byte("posts")
)

func init() {
	go func() {
		var err error
		db, err = bolt.Open("hentai.db", 0600, nil)
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		if err := db.Batch(func(tx *bolt.Tx) error {
			users, err := tx.CreateBucketIfNotExists(bktUsers)
			if err != nil {
				return err
			}

			if _, err = tx.CreateBucketIfNotExists(bktPosts); err != nil {
				return err
			}

			users.Tx().ForEach(func(name []byte, bkt *bolt.Bucket) error {
				bID, _ := strconv.Atoi(string(name))
				for _, admin := range cfg["admins"].([]interface{}) {
					if int(admin.(float64)) == bID {
						log.F("change %d role to admin", bID)
						return ChangeRoleBD(bID, "admin")
					} else {
						for _, patron := range cfg["patrons"].([]interface{}) {
							if int(patron.(float64)) == bID {
								log.F("change %d role to patron", bID)
								return ChangeRoleBD(bID, "patron")
							}
						}
					}
					log.Ln("skipped role bucket")
				}
				return nil
			})

			return nil
		}); err != nil {
			panic(err.Error())
		}

		select {}
	}()
}

func CreateUserBD(id int) error {
	return db.Batch(func(tx *bolt.Tx) error {
		bkt, err := tx.Bucket(bktUsers).CreateBucket([]byte(strconv.Itoa(id)))
		if err != nil {
			return err
		}

		for _, admin := range cfg["admins"].([]interface{}) {
			if id == int(admin.(float64)) {
				bkt.Put([]byte("role"), []byte("admin"))
			} else {
				for _, patron := range cfg["patrons"].([]interface{}) {
					if id == int(patron.(float64)) {
						bkt.Put([]byte("role"), []byte("patron"))
					} else {
						bkt.Put([]byte("role"), []byte("user"))
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
	if err := db.Batch(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		return bkt.Put([]byte("lang"), []byte(lang))
	}); err != nil {
		log.Ln(err.Error())
		CreateUserBD(id)
		return ChangeLangBD(id, lang)
	}
	return nil
}

func ChangeRoleBD(id int, role string) error {
	if err := db.Batch(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		return bkt.Put([]byte("role"), []byte(role))
	}); err != nil {
		log.Ln(err.Error())
		CreateUserBD(id)
		return ChangeRoleBD(id, role)
	}
	return nil
}

func ChangeFilterDB(id int, nsfw bool) error {
	if err := db.Batch(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		return bkt.Put([]byte("nsfw"), strconv.AppendBool(nil, nsfw))
	}); err != nil {
		log.Ln(err.Error())
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
		log.Ln(err.Error())
		CreateUserBD(id)
		return GetUserDB(id)
	}
	return &usr, nil
}

func AddHitsDB(id int) error {
	if err := db.Batch(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bktUsers).Bucket([]byte(strconv.Itoa(id)))
		if bkt == nil {
			return fmt.Errorf("bucket not exist")
		}

		hits, _ := strconv.Atoi(string(bkt.Get([]byte("hits"))))
		hits++
		log.F("%d hits to %d user", hits, id)
		return bkt.Put([]byte("hits"), []byte(strconv.Itoa(hits)))
	}); err != nil {
		log.Ln(err.Error())
		CreateUserBD(id)
		return AddHitsDB(id)
	}
	return nil
}
