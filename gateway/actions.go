package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type action struct {
	ID                    bson.ObjectId `bson:"_id,omitempty"`
	NameSpace, User, Code string
}

// GetAction gets an action
func GetAction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := ps.ByName("action")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	objectID := bson.ObjectIdHex(id)

	userAction := action{}
	err := db.C("actions").FindId(objectID).One(&userAction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userAction)
}

// PostAction stores an action to the database
func PostAction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10e6))
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	userAction := action{}

	err = json.Unmarshal(body, &userAction)
	if err != nil {
		w.WriteHeader(422)
		return
	}

	userAction.ID = bson.NewObjectId()
	err = db.C("actions").Insert(userAction)
	if err != nil {
		log.Printf("could not store action: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
