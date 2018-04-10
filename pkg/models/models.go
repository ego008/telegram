package models

import "time"

type (
	Model struct {
		ID        int       `gorm:"column:id;type:integer;primary_key;unique;not null"`
		CreatedAt time.Time `gorm:"column:created_at;not null"`
		UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	}

	User struct {
		Model
		Locale    string     `gorm:"column:locale;default:'en';size:2"`
		Ratings   Ratings    `gorm:"foreignkey:ID;association_foreignkey:ID"`
		BlackList []Tag      `gorm:"foreignkey:ID;association_foreignkey:ID"`
		WhiteList []Tag      `gorm:"foreignkey:ID;association_foreignkey:ID"`
		Resources []Resource `gorm:"foreignkey:ID;association_foreignkey:ID"`
		Roles     []Role     `gorm:"foreignkey:ID;association_foreignkey:ID"`
	}

	Ratings struct {
		Model
		Safe         bool `gorm:"column:safe;default:1"`
		Questionable bool `gorm:"column:questionable;default:0"`
		Exlplicit    bool `gorm:"column:explicit;default:0"`
	}

	Tag struct {
		Model
		Tag string `gorm:"column:tag;unique;not null"`
	}

	Types struct {
		Model
		Image     bool `gorm:"column:images;default:1"`
		Animation bool `gorm:"column:animations;default:1"`
		Video     bool `gorm:"column:videos;default:1"`
	}

	Resource struct {
		Model
		Name    string `gorm:"column:name;size:4;unique;not null"`
		Token   string `gorm:"column:token;unique"`
		Enabled bool   `gorm:"column:videos;default:0"`
	}

	Role struct {
		Model
		Role string `gorm:"column:role;unique;not null"`
	}

	Result struct {
		Source       string `json:"source"`
		Directory    string `json:"directory"`
		Hash         string `json:"hash"`
		Height       int    `json:"height"`
		ID           int    `json:"id"`
		Image        string `json:"image"`
		Change       int64  `json:"change"`
		Owner        string `json:"owner"`
		ParentID     int    `json:"parent_id"`
		Rating       string `json:"rating"`
		Sample       bool   `json:"sample"`
		SampleHeight int    `json:"sample_height"`
		SampleWidth  int    `json:"sample_width"`
		Score        int    `json:"score"`
		Tags         string `json:"tags"`
		Width        int    `json:"width"`
	}
)

func (User) TableName() string { return "users" }

func (Ratings) TableName() string { return "ratings" }
