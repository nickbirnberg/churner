package main

import (
	"log"
	"net/http"

	"github.com/nickbirnberg/churner/common"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func main() {
	// Create database connection
	log.Printf("Establishing MongoDB Connection")
	session, err := mgo.Dial(common.MustGetenv("MONGO_ADDR"))
	if err != nil {
		log.Fatalf("Failed to open MongoDB session: %v", err)
	}

	db = session.DB(common.MustGetenv("MONGO_DB_NAME"))

	defer func() {
		log.Println("Closing MongoDB Connection")
		session.Close()
	}()

	router := getRouter()

	log.Fatal(http.ListenAndServe(common.MustGetenv("HOST_ADDR"), router))
}
