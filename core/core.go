// Package core holds the business logic for yps
package core

import (
	"log"
	"sync"

	"github.com/gophergala/yps/queue"
)

const (
	// UserInputQueue represents the name of the user input queue
	UserInputQueue = `userInput`
)

func processMessage(task *queue.Message, mq *queue.Queue, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	// TODO check if message is for a playlist or video and send to the right mq
	log.Printf("Got task: %#q", (*task).Original())
	err = (*mq).Delete(task)

	return
}

// ProcessUserInput processes all the messages from the user input queue
func ProcessUserInput(mq *queue.Queue, resp chan<- error) {
	tasks, err := (*mq).Fetch(10)

	if err != nil {
		log.Printf("[error] Task failed: %#v", err)
		resp <- err
		return
	}

	wg := &sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task queue.Message) {
			if err := processMessage(&task, mq, wg); err != nil {
				log.Printf("[error] task failed to be processed task: %#v message: %q", task, err)
			}
		}(*task)

	}
	wg.Wait()

	resp <- nil
}
