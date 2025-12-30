package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"restapi/chapter_07/store_api/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type DBClient struct {
	db *gorm.DB
}

type UserResponse struct {
	User models.User `json:"user"`
	Data interface{} `json:"data"`
}

func (driver *DBClient) GetUsersByFirstName(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	name := r.FormValue("first_name")

	var query = "select * from user where data ->> 'first_name' = ?"
	driver.db.Raw(query, name).Scan(&users)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	respJSON, _ := json.Marshal(users)
	w.Write(respJSON)
}

func (driver *DBClient) GetUser(w http.ResponseWriter, r *http.Request) {
	var user = models.User{}
	vars := mux.Vars(r)

	driver.db.First(&user, vars["id"])
	var userData interface{}

	json.Unmarshal([]byte(user.Data), &userData)
	var response = UserResponse{User: user, Data: userData}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	respJson, _ := json.Marshal(response)
	w.Write(respJson)
}

func (driver *DBClient) PostUser(w http.ResponseWriter, r *http.Request) {
	var user = models.User{}
	postBody, _ := ioutil.ReadAll(r.Body)

	user.Data = string(postBody)
	driver.db.Save(&user)
	responseMap := map[string]interface{}{"id": user.ID}
	var err string = ""
	if err != "" {
		w.Write([]byte("yes"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(responseMap)
		w.Write(response)
	}
}

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	dbclient := &DBClient{db: db}
	defer models.CloseDB(db)

	r := mux.NewRouter()
	
	r.HandleFunc("/v1/user/{id:[a-zA-Z0-9]*}", dbclient.GetUser).Methods("GET")
	r.HandleFunc("/v1/user", dbclient.PostUser).Methods("POST")
	r.HandleFunc("/v1/user", dbclient.GetUsersByFirstName).Methods("GET")
	
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8004",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server starting on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}