package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func HealhcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, time.Now().String())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", HealhcheckHandler)
	srv := &http.Server{
		Handler: r,
		Addr: "0.0.0.0:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}