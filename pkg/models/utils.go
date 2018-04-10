package models

func (user *User) getTags(list []Tag) []string {
	tags := make([]string, len(list))
	for i := range list {
		tags[i] = list[i].Tag
	}
	return tags
}

func (user *User) BlackTags() []string {
	return user.getTags(user.BlackList)
}

func (user *User) WhiteTags() []string {
	return user.getTags(user.WhiteList)
}
