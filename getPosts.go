package main

import (
	"encoding/json"
	"log"
	"net/http"
	//"github.com/valyala/fasthttp" // Need to replace "net/http" later
)

// Universal(?) function obtain content
func getPosts(tags string, pid string) []Post {
	// JSON API with 50 results (Telegram limit)
	repository := config.Resource[resNum].Settings.URL + "index.php?page=dapi&s=post&q=index&json=1&limit=50"
	if tags != "" {
		repository += "&tags=" + tags // Insert tags
	}
	if pid != "" {
		repository += "&pid=" + pid // Insert result-page
	}
	resp, err := http.Get(repository) // Need to replace on "fasthttp" later :\
	if err != nil {
		log.Printf("Error in GET request: %s", err)
	}
	defer resp.Body.Close()
	var obj []Post
	json.NewDecoder(resp.Body).Decode(&obj)
	return obj
}
