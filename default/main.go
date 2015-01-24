package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"runtime"
	"path"
	"io/ioutil"
)

var (
	indexTemplate string
)

func init() {
	loadTemplates()

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexTemplate)
}
