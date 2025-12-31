package main

import (
	"context"
	"log"
	"os"

	"github.com/levigross/grequests/v2"
)

var GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
var requestOptions = &grequests.RequestOptions{Auth: []string{GITHUB_TOKEN, "x-oauth-basic"}}

type Repo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	FullName string `json:"full_name"`
	Forks int `json:"forks"`
	Private bool `json:"private"`
}

func getStatus(url string) *grequests.Response{
	resp, err := grequests.Get(context.Background(), url,
        grequests.UserAgent("MyAgent"))
	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}
	return resp
}

func main() {
	var repos []Repo
	var repoURL = "https://api.github.com/users/torvalds/repos"
	resp := getStatus(repoURL)
	resp.JSON(&repos)
	log.Println(repos)
}