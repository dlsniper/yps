// Package main holds the frontend logic for the user interaction
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"

	"appengine"
	"appengine/taskqueue"

	"github.com/gophergala/yps/core"
	ypshu "github.com/gophergala/yps/core/httputil"
	"github.com/gophergala/yps/provider/youtube"
	"github.com/gophergala/yps/queue/aetq"

	"github.com/gorilla/mux"
)

var (
	indexTemplate             string
	errInvalidMessageReceived = fmt.Errorf("invalid message received")
	errInternalServerError    = fmt.Errorf("internal server error")
)

func init() {
	loadTemplates()

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/addToQueue", addToQueue).Methods("POST")
	http.Handle("/", r)
}

func loadTemplates() (err error) {
	_, currentFilename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("Could not retrieve the caller for loading templates")
	}

	currentDir := path.Dir(currentFilename)
	var template []byte
	template, err = ioutil.ReadFile(path.Join(currentDir, "resources/index.html"))
	if err != nil {
		return
	}
	indexTemplate = string(template)

	return
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexTemplate)
}

func addToQueue(w http.ResponseWriter, r *http.Request) {
	yt := youtube.NewYoutube()

	if err := r.ParseForm(); err != nil {
		if !appengine.IsDevAppServer() {
			err = errInvalidMessageReceived
		}
		ypshu.WriteResponse(err, http.StatusBadRequest, r, w)
		return
	}

	url := r.PostForm.Get("url")

	if !yt.IsValidURL(url) {
		ypshu.WriteResponse(errInvalidMessageReceived, http.StatusBadRequest, r, w)
		return
	}

	payload, err := core.ProcessUserInput(url)

	if err != nil {
		if !appengine.IsDevAppServer() {
			err = errInvalidMessageReceived
		}

		ypshu.WriteResponse(err, http.StatusInternalServerError, r, w)
		return
	}

	msg := aetq.NewMessage(&taskqueue.Task{
		Payload: payload,
		Method:  "PULL",
	})

	c := appengine.NewContext(r)
	mq := aetq.NewQueue(c, core.UserInputQueue, 60)
	if err := mq.Add(&msg); err != nil {
		if !appengine.IsDevAppServer() {
			err = errInternalServerError
		}

		ypshu.WriteResponse(err, http.StatusInternalServerError, r, w)
		return
	}

	ypshu.WriteResponse(fmt.Sprintf("%q", "created"), http.StatusCreated, r, w)
}
