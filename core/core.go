// Package core holds the business logic for yps
package core

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gophergala/yps/queue"
)

type (
	userInput struct {
		URL string `json:"url"`
	}
)

const (
	// UserInputQueue represents the name of the user input queue
	UserInputQueue = `userInput`
)

func encodeUserInputTask(url string) ([]byte, error) {
	msg := &userInput{
		URL: url,
	}

	return json.Marshal(msg)
}

func decodeUserInputTask(msg string) (m userInput, err error) {
	err = json.Unmarshal([]byte(msg), &m)
	return
}

// ProcessUserInput transforms the input taken fro the user and returns it in the format needed
func ProcessUserInput(url string) ([]byte, error) {
	return encodeUserInputTask(url)
}

func processMessage(task *queue.Message, mq *queue.Queue, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	// TODO check if message is for a playlist or video and send to the right mq
	log.Printf("Got task: %#q", (*task).Original())
	err = (*mq).Delete(task)

	return
}

// ProcessUserInputTasks processes all the messages from the user input queue
func ProcessUserInputTasks(mq *queue.Queue, resp chan<- error) {
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
