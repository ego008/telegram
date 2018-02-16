package models

type User struct {
	Blacklist    []string
	ID           int
	Language     string
	Ratings      Ratings
	Resources    map[string]bool
	Roles        Roles
	ContentTypes ContentTypes
	Whitelist    []string
}
