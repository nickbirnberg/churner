package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nickbirnberg/churner/common"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	server = httptest.NewServer(getRouter())
	session, err := mgo.Dial(common.MustGetenv("MONGO_ADDR"))
	if err != nil {
		os.Exit(70741)
	}
	db = session.DB("test")

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
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("response status code '%v' unexpected", resp.StatusCode)
	}
}
