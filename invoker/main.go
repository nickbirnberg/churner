package main

import (
	"log"
	"net/http"

	"github.com/fsouza/go-dockerclient"
	"github.com/julienschmidt/httprouter"
	"github.com/nickbirnberg/churner/common"
	"gopkg.in/mgo.v2"
)

var db *mgo.Database
var dclient *docker.Client

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

	// Create Docker client
	log.Printf("Creating Docker client from ENV variables")
	dclient, err = docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	router := getRouter()

	log.Fatal(http.ListenAndServe(common.MustGetenv("HOST_ADDR"), router))
}

func getRouter() http.Handler {
	router := httprouter.New()

	router.POST("/invoke/:action", InvokeAction)

	return router
}
