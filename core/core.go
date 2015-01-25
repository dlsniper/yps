// Package core holds the business logic for yps
package core

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gophergala/yps/provider/youtube"
	"github.com/gophergala/yps/queue"
	"github.com/gophergala/yps/queue/aetq"
)

type (
	userInput struct {
		URL string `json:"url"`
	}
)

const (
	// UserInputQueue holds the name of the user input queue
	UserInputQueue = `userInput`
	// PlaylistQueue holds the name of the playlists queue
	PlaylistQueue = `playlist`
	// VideoQueue holds the name of the video queue
	VideoQueue = `video`
)

var (
	errInvalidYoutubeURL      = fmt.Errorf("invalid youtube url received")
	errTaskNotVideoOrPlaylist = fmt.Errorf("task was not a video or a playlist")
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

func processMessage(task *queue.Message, msgMq, playlistMq, videoMq *queue.Queue, wg *sync.WaitGroup) (err error) {
	defer func() {
		er := (*msgMq).Delete(task)
		if err == nil {
			err = er
		}
		wg.Done()
	}()

	log.Printf("Got task: %#q", (*task).Original())

	var msg userInput
	msg, err = decodeUserInputTask((*task).String())
	if err != nil {
		return
	}

	yt := youtube.NewYoutube()
	if !yt.IsValidURL(msg.URL) {
		return errInvalidYoutubeURL
	}

	message := aetq.NewMessage(msg.URL)

	if yt.IsPlaylist(msg.URL) {
		(*playlistMq).Add(&message)
	} else if yt.IsVideo(msg.URL) {
		(*videoMq).Add(&message)
	} else {
		err = errTaskNotVideoOrPlaylist
	}

	return
}

// ProcessUserInputTasks processes all the messages from the user input queue
func ProcessUserInputTasks(msgMq, playlistMq, videoMq *queue.Queue, resp chan<- error) {
	tasks, err := (*msgMq).Fetch(10)

	if err != nil {
		log.Printf("[error] Task failed: %#v", err)
		resp <- err
		return
	}

	wg := &sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task queue.Message) {
			if err := processMessage(&task, msgMq, playlistMq, videoMq, wg); err != nil {
				log.Printf("[error] task failed to be processed task: %#v message: %q", task, err)
			}
		}(*task)

	}
	wg.Wait()

	resp <- nil
}
