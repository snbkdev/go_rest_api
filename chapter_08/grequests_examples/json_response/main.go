package main

import (
	"context"
	"log"

	"github.com/levigross/grequests/v2"
)

func main() {
    resp, err := grequests.Get(context.Background(), "https://httpbin.org/get",
        grequests.UserAgent("MyAgent"))
    if err != nil {
        log.Fatal(err)
    }

    var data map[string]interface{}
	resp.JSON(&data)
	log.Println(data)
}