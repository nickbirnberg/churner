package main

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

func main() {
	// Create database connection
	log.Printf("Establishing MongoDB Connection")
	session, err := mgo.Dial(mustGetenv("MONGO_ADDR"))
	if err != nil {
		log.Fatalf("Failed to open MongoDB session: %v", err)
	}

	db = session.DB(mustGetenv("MONGO_DB_NAME"))

	defer func() {
		log.Println("Closing MongoDB Connection")
		session.Close()
	}()

	router := getRouter()

	log.Fatal(http.ListenAndServe(mustGetenv("HOST_ADDR"), router))
}

func mustGetenv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		log.Panicln("Missing env variable:", name)
	}
	return env
}
