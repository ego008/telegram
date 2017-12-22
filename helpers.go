package main

import "fmt"

func (usr *user) getRatingsStatus() (string, error) {
	T, err := langSwitch(usr.Language)
	if err != nil {
		return "", err
	}

	switch {
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe only
		return T("rating_safe"), nil
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// questionable only
		return T("rating_questionable"), nil
	case !usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// explicit only
		return T("rating_explicit"), nil
	case usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe + questionable
		return fmt.Sprintln(T("rating_safe"), "+", T("rating_questionable")), nil
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// safe + explicit
		return fmt.Sprintln(T("rating_safe"), "+", T("rating_explicit")), nil
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// questionable + explicit
		return fmt.Sprintln(T("rating_questionable"), "+", T("rating_explicit")), nil
	default:
		// all ratings enabled/diabled
		return T("rating_all"), nil
	}
}

func (usr *user) getRatingsFilter() string {
	switch {
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe only
		return "rating:safe"
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// questionable only
		return "rating:questionable"
	case !usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// explicit only
		return "rating:explicit"
	case usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		!usr.Ratings.Exlplicit:
		// safe + questionable
		return "-rating:explicit"
	case usr.Ratings.Safe &&
		!usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// safe + explicit
		return "-rating:questionable"
	case !usr.Ratings.Safe &&
		usr.Ratings.Questionable &&
		usr.Ratings.Exlplicit:
		// questionable + explicit
		return "-rating:safe"
	default:
		// all ratings enabled/diabled
		return ""
	}
}

func checkInterface(src interface{}) string {
	if src != nil {
		return src.(string)
	}
	return ""
}
