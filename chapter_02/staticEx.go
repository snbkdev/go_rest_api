package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main6() {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("/Users/asanbeksamudin/Documents/projects/mine/golang/building_rest_api/chapter_02/static"))
	log.Fatal(http.ListenAndServe(":8000", router))
}