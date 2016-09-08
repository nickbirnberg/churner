package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fsouza/go-dockerclient"
	"github.com/nickbirnberg/churner/common"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var server *httptest.Server

type responseError struct {
	Response, Error string
}

func TestMain(m *testing.M) {
	server = httptest.NewServer(getRouter())
	session, err := mgo.Dial(common.MustGetenv("MONGO_ADDR"))
	if err != nil {
		os.Exit(70741)
	}
	db = session.DB("test")

	dclient, err = docker.NewClientFromEnv()
	if err != nil {
		os.Exit(5000)
	}

	retCode := m.Run()

	session.Close()

	os.Exit(retCode)
}

func TestInvokeNonExistantAction(t *testing.T) {
	err := db.DropDatabase()
	if err != nil {
		t.Error(err)
	}

	resp, err := http.Post(server.URL+"/invoke/"+bson.NewObjectId().Hex(), "", new(bytes.Buffer))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("response status code '%v' unexpected", resp.StatusCode)
	}
}

func TestInvokePythonAction(t *testing.T) {
	err := db.DropDatabase()
	if err != nil {
		t.Error(err)
	}

	pythonAction := &common.Action{User: "testUser"}
	pythonAction.ID = bson.NewObjectId()
	pythonAction.Code = "def action_func(param): return param"
	err = db.C("actions").Insert(pythonAction)
	if err != nil {
		t.Errorf("failed to insert into db: %v", err)
	}

	var jsonString = []byte(`{"Param":"hello world"}`)
	resp, err := http.Post(server.URL+"/invoke/"+pythonAction.ID.Hex(), "application/json", bytes.NewBuffer(jsonString))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("response status code '%v' unexpected", resp.StatusCode)
	}

	var invokerResponse responseError
	json.NewDecoder(resp.Body).Decode(&invokerResponse)
	if invokerResponse.Response != "hello world" {
		t.Errorf("response body not expected: %v", invokerResponse.Response)
	}
}
