package db

import "github.com/HentaiDB/HentaiDBot/internal/errors"

func Close() {
	err := DB.Close()
	errors.Check(err)
}
