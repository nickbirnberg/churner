package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	server = httptest.NewServer(getRouter())
	session, err := mgo.Dial(mustGetenv("MONGO_ADDR"))
	if err != nil {
		os.Exit(70741)
	}
	db = session.DB("test")

	retCode := m.Run()

	server.Close()
	session.Close()

	os.Exit(retCode)
}

func TestGetEmptyCollection(t *testing.T) {
	err := db.DropDatabase()
	if err != nil {
		t.Errorf("could not drop DB: %v", err)
	}

	resp, err := http.Get(server.URL + "/api/v1/actions/" + bson.NewObjectId().Hex())
	if err != nil {
		t.Fatalf("could not get action: %v", err)
	}

	if resp.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		t.Fatalf("header content-type %v not expected", resp.Header.Get("Content-Type"))
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("response code '%v' not expected", resp.StatusCode)
	}
}

func TestGetAction(t *testing.T) {
	err := db.DropDatabase()
	if err != nil {
		t.Errorf("could not drop DB: %v", err)
	}

	someAction := &action{User: "testUser"}
	someAction.ID = bson.NewObjectId()
	err = db.C("actions").Insert(someAction)
	if err != nil {
		t.Errorf("failed to insert into db: %v", err)
	}

	someActionjson, err := json.Marshal(someAction)
	if err != nil {
		t.Errorf("could not marshal: %v", err)
	}

	resp, err := http.Get(server.URL + "/api/v1/actions/" + someAction.ID.Hex())
	if err != nil {
		t.Errorf("error getting action: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("response code '%v' not expected", resp.Status)
	}

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if bytes.Equal(jsonResp, someActionjson) {
		t.Errorf("response %s did not match %s", jsonResp, someActionjson)
	}
}

func TestPostAction(t *testing.T) {
	err := db.DropDatabase()
	if err != nil {
		t.Errorf("could not drop DB: %v", err)
	}

	someAction := &action{User: "testUser"}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(*someAction)

	if err != nil {
		t.Errorf("could not marshal: %v", err)
	}

	resp, err := http.Post(server.URL+"/api/v1/actions", "application/json; charset=utf-8", b)
	if err != nil {
		t.Errorf("could not post action: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("response code '%v' not expected", resp.Status)
	}

	returnedAction := &action{}
	err = db.C("actions").Find(nil).One(&returnedAction)
	if err != nil {
		t.Error(err)
	}

	if returnedAction.User != someAction.User {
		t.Errorf("returned struct \n%+v does not match inserted struct \n%+v", returnedAction, someAction)
	}
}
