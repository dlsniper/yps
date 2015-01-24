package main

import (
	"fmt"
	"log"
	"net/http"

	"appengine"
	"appengine/taskqueue"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/processUserInput", queueHandler).Methods("GET")
	http.Handle("/", r)
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	tasks, err := taskqueue.Lease(c, 10, "userInput", 60)

	if err != nil {
		log.Printf("[error] Task failed: %#v", err)
		return
	}

	for _, task := range tasks {
		log.Printf("Got task: %#v", task)
		taskqueue.Delete(c, task, "userInput")
	}

	fmt.Fprint(w, "it worked")
}
