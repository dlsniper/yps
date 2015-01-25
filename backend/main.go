// Package backend implements the main for the backend requests
package backend

import (
	"fmt"
	"log"
	"net/http"

	"appengine"

	"github.com/gophergala/yps/queue/aetq"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/processUserInput", queueHandler).Methods("GET")
	http.Handle("/", r)
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	queue := aetq.NewQueue(c, "userInput", 60)

	tasks, err := queue.Fetch(60)

	if err != nil {
		log.Printf("[error] Task failed: %#v", err)
		return
	}

	for _, task := range tasks {
		log.Printf("Got task: %#q", (*task).Original())
		queue.Delete(task)
	}

	fmt.Fprint(w, "it worked")
}
