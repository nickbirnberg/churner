package main

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
	"github.com/nickbirnberg/churner/common"
)

// InvokeAction invokes an action.
func InvokeAction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := ps.ByName("action")
	objectID := bson.ObjectIdHex(id)

	userAction := common.Action{}
	err := db.C("actions").FindId(objectID).One(&userAction)
	if err == mgo.ErrNotFound {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("query for %v unsuccessful :%v\n", objectID, err)
	}
}
