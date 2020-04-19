package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo configuration
const (
	hosts      = "localhost:27017"
	database   = "db"
	username   = ""
	password   = ""
	collection = "jobs"
)

// Job struct saved in the DB
type Job struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	Salary      string `json:"salary"`
}

type MongoStore struct {
	session *mgo.Session
}

var mongoStore = MongoStore{}

func main() {

	//Create MongoDB session
	session := initialiseMongo()
	mongoStore.session = session

	// define endpoints
	//router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/jobs", jobsGetHandler).Methods("GET")
	//router.HandleFunc("/jobs", jobsPostHandler).Methods("POST")
	http.HandleFunc("/jobs", jobsHandler)
	log.Fatal(http.ListenAndServe(":9090", nil))

}
func jobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		jobsPostHandler(w, r)
	} else {
		jobsGetHandler(w, r)
	}

}

// Initialize the MongoDB Session with 60 second timeout
func initialiseMongo() (session *mgo.Session) {

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	return

}

// http get route handler
func jobsGetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	col := mongoStore.session.DB(database).C(collection)

	results := []Job{}
	col.Find(bson.M{"title": bson.RegEx{"", ""}}).All(&results)
	jsonString, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(jsonString))

}

// http post route handler
func jobsPostHandler(w http.ResponseWriter, r *http.Request) {

	col := mongoStore.session.DB(database).C(collection)

	//Retrieve body from http request
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	//Save data into Job struct
	var _job Job
	err = json.Unmarshal(b, &_job)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Insert job into MongoDB
	err = col.Insert(_job)
	if err != nil {
		panic(err)
	}

	//Convert job struct into json
	jsonString, err := json.Marshal(_job)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Set content-type http header
	w.Header().Set("content-type", "application/json")

	//Send back data as response
	w.Write(jsonString)

}
