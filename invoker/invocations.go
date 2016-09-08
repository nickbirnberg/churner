package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/fsouza/go-dockerclient"
	"github.com/julienschmidt/httprouter"
	"github.com/nickbirnberg/churner/common"
)

// InvokeAction invokes an action.
func InvokeAction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	defer r.Body.Close()

	// Get Code from DB
	actionCode, err := codeFromDB(ps.ByName("action"))
	if err == mgo.ErrNotFound {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get Params from request
	params, err := getParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create and run container
	containerID, err := createRunContainer()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	del := func() error {
		return dclient.RemoveContainer(docker.RemoveContainerOptions{ID: containerID, Force: true})
	}
	defer del()
	time.Sleep(300 * time.Millisecond) // todo: figure out better way to wait for container to start up

	// Post Code to container
	containerResponse, err := postAction(actionCode, params)

	w.Write(containerResponse)
}

func getParam(r *http.Request) (string, error) {
	var actionParams struct {
		Param string
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10e6))
	if err != nil {
		return "", fmt.Errorf("couldn't read request body:%v", err)
	}
	err = json.Unmarshal(body, &actionParams)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling user param:%v", err)
	}
	return actionParams.Param, nil
}

func createRunContainer() (string, error) {
	containerOptions := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: "nbirnberg/churner-py3",
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{
				"5000/tcp": []docker.PortBinding{docker.PortBinding{HostPort: "8091"}},
			},
		},
	}
	container, err := dclient.CreateContainer(containerOptions)
	if err != nil {
		return "", fmt.Errorf("failed to create container:%v", err)
	}
	err = dclient.StartContainer(container.ID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to start container:%v", err)
	}

	return container.ID, nil
}

func postAction(code, params string) ([]byte, error) {
	containerCode := struct {
		Code, Param string
	}{
		code, params,
	}

	jsonBuffer := new(bytes.Buffer)
	err := json.NewEncoder(jsonBuffer).Encode(containerCode)
	if err != nil {
		return nil, fmt.Errorf("could not encode json to send to container:%v", err)
	}

	resp, err := http.Post("http://127.0.0.1"+":8091"+"/run", "application/json", jsonBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to POST to container:%v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can not read response from container:%v", err)
	}
	return body, nil
}

func codeFromDB(id string) (string, error) {
	objectID := bson.ObjectIdHex(id)
	userAction := common.Action{}
	err := db.C("actions").FindId(objectID).One(&userAction)
	if err == mgo.ErrNotFound {
		return "", err
	}
	if err != nil {
		log.Printf("query for %v unsuccessful :%v", objectID, err)
		return "", err
	}

	return userAction.Code, nil
}
