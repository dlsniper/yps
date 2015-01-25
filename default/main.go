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

	"github.com/gophergala/yps/provider/youtube"
	"github.com/gophergala/yps/queue/aetq"

	"github.com/gorilla/mux"
)

var (
	indexTemplate string
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

	if r.ParseForm() != nil {
		writeResponse(fmt.Errorf("invalid message received"), http.StatusBadRequest, r, w)
		return
	}

	url := r.PostForm.Get("url")

	if !yt.IsValidURL(url) {
		writeResponse(fmt.Errorf("invalid message received"), http.StatusBadRequest, r, w)
		return
	}

	msg := aetq.NewMessage(&taskqueue.Task{
		Payload: []byte(url),
		Method:  "PULL",
	})

	c := appengine.NewContext(r)
	queue := aetq.NewQueue(c, "userInput", 60)
	if err := queue.Add(&msg); err != nil {
		if appengine.IsDevAppServer() {
			err = fmt.Errorf("%q", err)
		}
		writeResponse(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(fmt.Sprintf("%q", "created"), http.StatusCreated, r, w)
}

func writeResponse(response interface{}, code int, r *http.Request, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	fmt.Fprintf(w, "%d %q", code, response)
}
