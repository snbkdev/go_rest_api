package main

import (
	"log"
	"restapi/chapter_07/models"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Println(db)
	}
}