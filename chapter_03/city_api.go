package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type city2 struct {
	Name string
	Area uint64
}

func filterContentType(handler http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("currently in the check content type middleware")
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - Unsupported Media Type. Please send JSON"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func setServerTimeCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)

		cookie := http.Cookie{Name: "Server-Time(UTC)", Value: strconv.FormatInt(time.Now().Unix(), 10)}
		http.SetCookie(w, &cookie)
		log.Println("Currently in the set server time middleware")
	})
}

func mainLogica2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tempCity city2
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&tempCity)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		log.Printf("Got %s city with area of %d sq miles!\n", tempCity.Name, tempCity.Area)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 - Created"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}
}

func main6() {
	mainLogicHandler := http.HandlerFunc(mainLogica2)
	http.Handle("/city", filterContentType(setServerTimeCookie(mainLogicHandler)))
	http.ListenAndServe(":8000", nil)
}