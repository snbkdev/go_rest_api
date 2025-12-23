package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func main4() {
	newMux := http.NewServeMux()

	newMux.HandleFunc("/randomfloat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rand.Float64())
	})

	newMux.HandleFunc("/randomint", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, rand.Int())
	})

	http.ListenAndServe(":8000", newMux)
}