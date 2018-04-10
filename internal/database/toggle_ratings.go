package database

import (
	"github.com/HentaiDB/HentaiDBot/pkg/models"
	tg "github.com/toby3d/telegram"
)

func (db *DataBase) toggleRating(usr *tg.User, rating string) error {
	user, err := db.GetUser(usr)
	if err != nil {
		return err
	}

	model := db.Model(user.Ratings)
	switch {
	case rating == models.RatingSafe:
		return model.Set(models.RatingSafe, !user.Ratings.Safe).Error
	case rating == models.RatingQuestionable:
		return model.Set(models.RatingSafe, !user.Ratings.Safe).Error
	case rating == models.RatingExplicit:
		return model.Set(models.RatingSafe, !user.Ratings.Safe).Error
	default:
		return nil
	}
}

func (db *DataBase) ToggleRatingSafe(usr *tg.User) error {
	return db.toggleRating(usr, models.RatingSafe)
}

func (db *DataBase) ToggleRatingQuestionable(usr *tg.User) error {
	return db.toggleRating(usr, models.RatingQuestionable)
}

func (db *DataBase) ToggleRatingExplicit(usr *tg.User) error {
	return db.toggleRating(usr, models.RatingExplicit)
}
