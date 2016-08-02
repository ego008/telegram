package main

const (
	startMessage              = "Hi, %s!\nI'm the official bot of @HentaiDB project.\n\nYou can search images and GIF's in Gelbooru through me. Just type `@HentaiDBot hatsune_miku` in any chat and select result which you want to send. 😏\n\nTell me /cheatsheet here to get a cheat sheet of advanced search.\n\nI become better due to support of volunteers, to which you can join in [GitHub-repository](https://github.com/toby3d/HentaiDBot). ☺️"
	helpMessage               = "/cheatsheet - Get cheat sheet of advanced search in Gelbooru."
	cheatSheetMessage         = "*tag1 tag2*\nSearch for posts that have `tag1` and `tag2`.\n\n*~tag1 ~tag2*\nSearch for posts that have `tag1` or `tag2`. (Currently does not work)\n\n*night~*\nFuzzy search for the tag `night`. This will return results such as `night fight bright` and so on according to the [Levenshtein distance](http://en.wikipedia.org/wiki/Levenshtein_distance).\n\n*-tag1*\nSearch for posts that don't have `tag1`.\n\n*ta∗1*\nSearch for posts with tags that starts with `ta` and ends with `1`.\n\n*user:bob*\nSearch for posts uploaded by the user `Bob`.\n\n*md5:foo*\nSearch for posts with the MD5 hash `foo`.\n\n*md5:foo∗*\nSearch for posts whose MD5 starts with the MD5 hash `foo`.\n\n*rating:questionable*\nSearch for posts that are rated `questionable`.\n\n*-rating:questionable*\nSearch for posts that are not rated `questionable`.\n\n*parent:1234*\nSearch for posts that have `1234` as a parent (and include post `1234`).\n\n*rating:questionable rating:safe*\nIn general, combining the same metatags (the ones that have colons in them) will not work.\n\n*rating:questionable parent:100*\nYou can combine different metatags, however.\n\n*width:>=1000 height:>1000*\nFind images with a width greater than or equal to `1000` and a height greater than `1000`.\n\n*score:>=10*\nFind images with a score greater than or equal to `10`. This value is updated once daily at 12AM CST.\n\n*sort:updated:desc*\nSort posts by their most recently updated order.\n\n*Other sortable types:*\n• `id`\n• `score`\n• `rating`\n• `user`\n• `height`\n• `width`\n• `parent`\n• `source`\n• `updated`\nCan be sorted by both `asc` or `desc`."
	feedbackMessage           = "*New feedback message!*\n%s\n- @%s"
	feedbackQuestion          = "I really want to help you. 😔\nWhat is wrong with me?"
	feedbackAnswer            = "Thanks for you feedback, %s! ☺️\nSempai sure to notice you!"
	feedbackEmpty             = "I can't bother Senpai without reason. 😕\nDescribe your request or suggestion a different way."
	noInlineResultTitle       = "Nobody here but us chickens!"
	noInlineResultDescription = "Try search a different combination of tags."
	noInlineResultMessage     = "Sumimasen, but, unfortunately I could not find desired content. 😓\nBut perhaps this it already present in @HentaiDB channel or his official group."
)
