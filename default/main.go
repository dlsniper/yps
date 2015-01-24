package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"runtime"

	"appengine"
	"appengine/taskqueue"

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
	r.HandleFunc("/addPlaylist", addPlaylistHandler).Methods("GET")
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

func addPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	queue := aetq.NewQueue(c, "userInput", 60)

	msg := aetq.NewMessage(&taskqueue.Task{
		Payload: []byte("hello world"),
		Method:  "PULL",
	})

	_ = queue.Add(&msg)
}
