package db

import (
	"github.com/HentaiDB/HentaiDBot/internal/models"
	"github.com/tidwall/buntdb"
)

func GetUserElseAdd(id int, lang string) (usr *models.User, err error) {
	usr, err = User(id)
	if err == buntdb.ErrNotFound {
		usr, err = AddUser(id, lang)
	}

	return
}
