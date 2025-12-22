package main

import (
	"fmt"
	"html"
	"net/http"
	romannumerals "restapi/chapter_01/romanNumerals"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPAthElements := strings.Split(r.URL.Path, "/")
		if urlPAthElements[1] == "roman_number" {
			number, _ := strconv.Atoi(strings.TrimSpace(urlPAthElements[2]))
			if number == 0 || number > 10 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Not found"))
			} else {
				fmt.Fprintf(w, "%q", html.EscapeString(romannumerals.Numerals[number]))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad Request"))
		}
	})

	s := &http.Server{
		Addr: ":7050",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}