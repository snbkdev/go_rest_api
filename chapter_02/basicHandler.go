package main

import (
	"io"
	"log"
	"net/http"
)

func MyServer(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, REST api !!!\n")
}

func main2() {
	http.HandleFunc("/hello", MyServer)
	log.Fatal(http.ListenAndServe(":7001", nil))
}