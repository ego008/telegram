package database

import (
	// "github.com/HentaiDB/HentaiDBot/pkg/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DataBase struct{ *gorm.DB }

var DB *DataBase

func Open(path string) (*DataBase, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.SingularTable(false)
	db.LogMode(true)

	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	// db.CreateTable(&models.Post{})

	return &DataBase{db}, nil
}
